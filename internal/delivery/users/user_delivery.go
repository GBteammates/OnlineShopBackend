package delivery

import (
	"OnlineShopBackend/internal/delivery/helper"
	"OnlineShopBackend/internal/delivery/users/user"
	"OnlineShopBackend/internal/delivery/users/user/googleOauth2"
	"OnlineShopBackend/internal/delivery/users/user/jwtauth"
	"OnlineShopBackend/internal/delivery/users/user/password"
	"OnlineShopBackend/internal/models"
	usecase "OnlineShopBackend/internal/usecase/interfaces"
	"errors"
	"fmt"
	"net/http"

	"github.com/dghubble/gologin/v2"
	gg "github.com/dghubble/gologin/v2/google"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"golang.org/x/oauth2"
	og2 "golang.org/x/oauth2/google"
)

type UserDelivery struct {
	userUsecase usecase.IUserUsecase
	cartUsecase usecase.ICartUsecase
	logger      *zap.SugaredLogger
}

func NewUserDelivery(userUsecase usecase.IUserUsecase, cartUsecase usecase.ICartUsecase, logger *zap.SugaredLogger) *UserDelivery {
	return &UserDelivery{
		userUsecase: userUsecase,
		cartUsecase: cartUsecase,
		logger:      logger,
	}
}

const (
	authorizationHeader = "Authorization"
)

// CreateUser create a new user
//
//	@Summary		Create a new user
//	@Description	Method provides to create a user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body	user.CreateUserData	true	"User data"
//	@Success		201
//	@Failure		400	{object}	ErrorResponse
//	@Failure		100	{object}	ErrorResponse
//	@Router			/user/create [post]
func (delivery *UserDelivery) CreateUser(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery CreateUser()")
	ctx := c.Request.Context()
	//var newUser *models.User
	var newUser *user.CreateUserData
	if err := c.ShouldBindJSON(&newUser); err != nil {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}

	newUserRights, err := delivery.userUsecase.GetRightsId(ctx, "Customer")
	if err != nil {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}
	modelsUser := models.User{
		Firstname: newUser.Firstname,
		Lastname:  newUser.Lastname,
		Password:  newUser.Password,
		Email:     newUser.Email,
		Address: models.UserAddress{
			Zipcode: newUser.Address.Zipcode,
			Country: newUser.Address.Country,
			City:    newUser.Address.City,
			Street:  newUser.Address.Street,
		},
		Rights: *newUserRights,
	}

	// Check user in database
	if existedUser, err := delivery.userUsecase.GetUserByEmail(ctx, newUser.Email); err == nil && existedUser.Id != uuid.Nil {
		delivery.logger.Error("user is already exists")
		helper.SetError(c, http.StatusContinue, fmt.Errorf("user is already exists"))
		return
	}

	// Password validation check
	if err := password.ValidationCheck(modelsUser); err != nil {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}

	hashPassword := password.GeneratePasswordHash(newUser.Password)
	newUser.Password = hashPassword

	// Create a user
	id, err := delivery.userUsecase.CreateUser(ctx, &modelsUser)
	if err != nil {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Infof("User with id: %v create success", id)

	c.JSON(http.StatusCreated, gin.H{"message": "success: user was created"})

}

// LoginUser login user
//
//	@Summary		Login user
//	@Description	Method provides to login a user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		password.Credentials	true	"Login"
//	@Success		200		{object}	user.LoginResponseData
//	@Failure		404		"Bad Request"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/user/login [post]
func (delivery *UserDelivery) LoginUser(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery LoginUser()")
	var userCredentials *user.Credentials

	if err := c.ShouldBindJSON(&userCredentials); err != nil {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()
	userExist, err := delivery.userUsecase.GetUserByEmail(ctx, userCredentials.Email)
	if err != nil || userExist.Password != password.GeneratePasswordHash(userCredentials.Password) {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusUnauthorized, fmt.Errorf("incorrect email or password"))
		return
	}

	if userExist.Email == "" {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusUnauthorized, err)
		return
	}

	cartExist, err := delivery.cartUsecase.GetCartByUserId(ctx, userExist.Id)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusContinue, err)
	}

	var cartId uuid.UUID
	if cartExist != nil {
		cartId = cartExist.Id
	} else {
		cartId, err = delivery.cartUsecase.CreateCart(ctx, userExist.Id)
		if err != nil {
			delivery.logger.Error(err.Error())
			helper.SetError(c, http.StatusNotFound, err)
		}
	}

	token, err := jwtauth.CreateSessionJWT(ctx, userExist)
	if err != nil {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusNotFound, err)
		return
	}

	res := user.LoginResponseData{
		CartId: cartId,
		Token: jwtauth.Token{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}

	c.JSON(http.StatusOK, res)
}

// GetUserProfile get user profile
//
//	@Summary		User profile
//	@Description	Method provides to get profile info
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth || firebase
//	@Success		201	{object}	user.CreateUserData
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Router			/user/profile [get]
func (delivery *UserDelivery) GetUserProfile(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UserProfile()")
	userCr, ok := c.MustGet("claims").(*jwtauth.Payload)
	if !ok {
		delivery.logger.Error("claims error")
		helper.SetError(c, http.StatusNotFound, nil)
		return
	}
	userData, err := delivery.userUsecase.GetUserByEmail(c.Request.Context(), userCr.Email)
	if err != nil {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusNotFound, err)
		return
	}

	password.SanitizePassword(userData)
	c.JSON(http.StatusCreated, userData)
}

// UpdateUserData user profile update
//
//	@Summary		User profile update
//	@Description	Method provides to update profile info
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body	user.CreateUserData	true	"New user data"
//	@Security		ApiKeyAuth || firebase
//	@Success		201	{object}	user.CreateUserData
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Router			/user/profile/edit [put]
func (delivery *UserDelivery) UpdateUserData(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UserProfileUpdate()")

	userCr, ok := c.MustGet("claims").(*jwtauth.Payload)
	if !ok {
		delivery.logger.Error("claims error")
		helper.SetError(c, http.StatusNotFound, nil)
		return
	}

	var newUserData *user.CreateUserData
	if err := c.ShouldBindJSON(&newUserData); err != nil {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}

	

	modelsUser := models.User{
		Id: newUserData.Id,
		Firstname: newUserData.Firstname,
		Lastname: newUserData.Lastname,
		Password: newUserData.Password,
		Address: models.UserAddress{
			Zipcode: newUserData.Address.Zipcode,
			Country: newUserData.Address.Country,
			City: newUserData.Address.City,
			Street: newUserData.Address.Street,
		},
	}

	userUpdated, err := delivery.userUsecase.UpdateUserData(c.Request.Context(), &modelsUser)
	if err != nil {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusNotFound, err)
		return
	}

	userUpdated.Email = userCr.Email
	userUpdated.Id = userCr.UserId

	c.JSON(http.StatusCreated, userUpdated)
}

// LoginUserGoogle Login Google
//
//	@Summary		Login with Google oauth2
//	@Description	Method provides to log in with Google
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		500
//	@Router			/user/login/google [get]
func (delivery *UserDelivery) LoginUserGoogle(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery LoginUserGoogle()")
	cfg, err := googleOauth2.NewUserConfig()
	if err != nil {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusNotFound, err)
		return
	}

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     og2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}

	stateConfig := gologin.DefaultCookieConfig
	gg.StateHandler(stateConfig, gg.LoginHandler(oauth2Config, nil)).ServeHTTP(c.Writer, c.Request)

}

func (delivery *UserDelivery) CallbackGoogle(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery CallbackGoogle()")
	cfg, err := googleOauth2.NewUserConfig()
	if err != nil {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusNotFound, err)
		return
	}

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     og2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}

	stateConfig := gologin.DefaultCookieConfig
	gg.StateHandler(stateConfig, gg.CallbackHandler(oauth2Config, delivery.success(c), delivery.failure(c))).ServeHTTP(c.Writer, c.Request)
}

func (delivery *UserDelivery) success(c *gin.Context) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		delivery.logger.Debug("Enter in delivery success()")
		ctx := req.Context()
		googleUser, err := gg.UserFromContext(ctx)
		if err != nil {
			delivery.logger.Error(err.Error())
			helper.SetError(c, http.StatusNotFound, err)
			return
		}
		if googleUser.Email == "" {
			delivery.logger.Error(err.Error())
			helper.SetError(c, http.StatusNotFound, err)
			return
		}
		u, err := delivery.userUsecase.GetUserByEmail(req.Context(), googleUser.Email)
		if err != nil {
			var newUser = models.User{
				Firstname: googleUser.GivenName,
				Lastname:  googleUser.FamilyName,
				Email:     googleUser.Email,
			}
			_, err := delivery.userUsecase.CreateUser(req.Context(), &newUser)
			if err != nil {
				delivery.logger.Error(err.Error())
				helper.SetError(c, http.StatusInternalServerError, err)
				return
			}
			u, err = delivery.userUsecase.GetUserByEmail(ctx, newUser.Email)
			if err != nil {
				delivery.logger.Error(err.Error())
				helper.SetError(c, http.StatusInternalServerError, err)
				return
			}
		}
		token, err := jwtauth.CreateSessionJWT(c.Request.Context(), u)
		if err != nil {
			delivery.logger.Error(err.Error())
			helper.SetError(c, http.StatusNotFound, err)
			return
		}

		c.JSON(http.StatusOK, token)

	}
	return fn
}

func (delivery *UserDelivery) failure(c *gin.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		delivery.logger.Error("google callback failure")
		c.JSON(http.StatusNotFound, gin.H{"error": "google callback failure"})
	}
}

// LogoutUser logout
//
//	@Summary		Logout
//	@Description	Method provides to log out
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		404	"Bad Request"
//	@Router			/user/logout [get]
func (delivery *UserDelivery) LogoutUser(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery LogoutUser()")
	c.Set(authorizationHeader, "")
	c.JSON(http.StatusOK, gin.H{"message": "you have been successfully logged out"})
}

// ChangeRole Change User Role
//
//	@Summary		Change User Role
//	@Description	Method provides to Change User Role
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		500
//	@Router			/user/callbackGoogle [put]
func (delivery *UserDelivery) ChangeRole(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery ChangeRights()")
	userCr, ok := c.MustGet("claims").(*jwtauth.Payload)
	if !ok {
		delivery.logger.Error("claims error")
		helper.SetError(c, http.StatusNotFound, nil)
		return
	}

	var newInfoUser *models.User
	if err := c.ShouldBindJSON(&newInfoUser); err != nil {
		delivery.logger.Error("claims error")
		helper.SetError(c, http.StatusNotFound, nil)
		return
	}

	if userCr.Email == newInfoUser.Email {
		c.JSON(http.StatusUnavailableForLegalReasons, gin.H{"error": "change your own rights is prohibited"})
		return
	}
	roleId, _ := delivery.userUsecase.GetRightsId(c.Request.Context(), newInfoUser.Rights.Name)

	err := delivery.userUsecase.UpdateUserRole(c.Request.Context(), roleId.Id, newInfoUser.Email)
	if err != nil {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "new user role"})

}

// RolesList List of Rights
//
//	@Summary		Change User Role
//	@Description	Method provides to Change User Role
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		500	{object}	ErrorResponse
//	@Failure		401	{object}	ErrorResponse
//	@Router			/user/rights/list [get]
func (delivery *UserDelivery) RolesList(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery RolesList()")
	userCr := c.MustGet("claims").(*jwtauth.Payload)
	if userCr.Role != "Admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not permitted"})
		return
	}
	roles, err := delivery.userUsecase.GetRightsList(c.Request.Context())
	if err != nil {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, roles)
}

// CreateRights
//
//	@Summary		Method provides to create rights
//	@Description	Method provides to create rights.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			rights	body		user.ShortRights	true	"Data for creating rights"
//	@Success		201		{object}	user.RightsId
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/user/createRights/ [post]
func (delivery *UserDelivery) CreateRights(c *gin.Context) {
	delivery.logger.Debugf("Enter in delivery CreateRights()")

	var createdRights user.ShortRights
	if err := c.ShouldBindJSON(&createdRights); err != nil {
		delivery.logger.Errorf("error on bind json from request: %v", err)
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}
	if createdRights.Name == "" {
		err := fmt.Errorf("empty name is not correct")
		if err != nil {
			delivery.logger.Error(err.Error())
			helper.SetError(c, http.StatusBadRequest, err)
			return
		}
	}
	ctx := c.Request.Context()
	id, err := delivery.userUsecase.CreateRights(ctx, &models.Rights{
		Name:  createdRights.Name,
		Rules: createdRights.Rules,
	})
	if err != nil {
		delivery.logger.Error(err.Error())
		helper.SetError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, user.RightsId{Value: id.String()})
}
