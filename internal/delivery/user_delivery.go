/*
 * Backend for Online Shop
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package delivery

import (
	"OnlineShopBackend/internal/delivery/user"
	"OnlineShopBackend/internal/delivery/user/userconfig"
	"OnlineShopBackend/internal/models"
	"net/http"

	"github.com/dghubble/gologin/v2"
	gg "github.com/dghubble/gologin/v2/google"
	"golang.org/x/oauth2/yandex"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"golang.org/x/oauth2"
	og2 "golang.org/x/oauth2/google"
	//"golang.org/x/oauth2/yandex"
)

func (delivery *Delivery) CreateUser(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery CreateUser()")
	ctx := c.Request.Context()
	var newUser *models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check user in database
	if existedUser, err := delivery.userUsecase.GetUserByEmail(ctx, newUser.Email); err == nil && existedUser.ID != uuid.Nil {
		c.JSON(http.StatusContinue, gin.H{"error": err.Error()})
		return
	}

	// Password validation check
	if err := newUser.ValidationCheck(); err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	hashPassword, err := newUser.GeneratePasswordHash()
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	newUser.Password = hashPassword

	// Create a user
	createdUser, err := delivery.userUsecase.CreateUser(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	delivery.logger.Info("success: user was created")

	token, err := delivery.userUsecase.CreateSessionJWT(c.Request.Context(), createdUser, delivery.secretKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token.AccessToken, 3600*24*30, "", "", false, true)

	//user.IssueSession(delivery.logger, createdUser.ID.String())

	c.JSON(http.StatusOK, token)
}

func (delivery *Delivery) LoginUser(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery LoginUser()")
	var userCredentials user.Credentials
	if err := c.ShouldBindJSON(&userCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userExist, err := delivery.userUsecase.GetUserByEmail(c.Request.Context(), userCredentials.Email) //TODO check password
	if err != nil || !userExist.CheckPasswordHash(userCredentials.Password) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if userExist.Email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := delivery.userUsecase.CreateSessionJWT(c.Request.Context(), userExist, delivery.secretKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token.AccessToken, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

func (delivery *Delivery) UserProfile(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UserProfile()")
	header := c.GetHeader("Authorization")
	userCr, err := delivery.userUsecase.UserIdentity(header)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	userData, err := delivery.userUsecase.GetUserByEmail(c.Request.Context(), userCr.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	userProfile := models.User{
		Email:     userData.Email,
		Firstname: userData.Firstname,
		Lastname:  userData.Lastname,
		Address: models.UserAddress{
			Zipcode: userData.Address.Zipcode,
			Country: userData.Address.Country,
			City:    userData.Address.City,
			Street:  userData.Address.Street,
		},
		Rights: models.Rights{
			ID:    userData.Rights.ID,
			Name:  userData.Rights.Name,
			Rules: userData.Rights.Rules,
		},
	}
	c.JSON(http.StatusCreated, userProfile)
}

func (delivery *Delivery) UserProfileUpdate(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UserProfileUpdate()")
	header := c.GetHeader("Authorization")
	userCr, err := delivery.userUsecase.UserIdentity(header)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newInfoUser *models.User
	if err = c.ShouldBindJSON(&newInfoUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser := models.User{
		ID:        userCr.UserId,
		Firstname: newInfoUser.Firstname,
		Lastname:  newInfoUser.Lastname,
		Address: models.UserAddress{
			Zipcode: newInfoUser.Address.Zipcode,
			Country: newInfoUser.Address.Country,
			City:    newInfoUser.Address.City,
			Street:  newInfoUser.Address.Street,
		},
	}

	userUpdated, err := delivery.userUsecase.UpdateUserData(c.Request.Context(), &updatedUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userUpdated.Email = userCr.Email

	c.JSON(http.StatusCreated, userUpdated)
}

func (delivery *Delivery) TokenUpdate(c *gin.Context) {

}

func (delivery *Delivery) LoginUserGoogle(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery LoginUserGoogle()")
	cfg, err := userconfig.NewUserConfig()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

func (delivery *Delivery) CallbackGoogle(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery CallbackGoogle()")
	cfg, err := userconfig.NewUserConfig()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     og2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}

	stateConfig := gologin.DefaultCookieConfig
	gg.StateHandler(stateConfig, gg.CallbackHandler(oauth2Config, delivery.success(c), failure(c))).ServeHTTP(c.Writer, c.Request)
}

func (delivery *Delivery) success(c *gin.Context) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		delivery.logger.Debug("Enter in delivery success()")
		ctx := req.Context()
		googleUser, err := gg.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if googleUser.Email == "" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u, err := delivery.userUsecase.GetUserByEmail(req.Context(), googleUser.Email)
		if err != nil {
			var NewUserCred = models.User{
				Firstname: googleUser.GivenName,
				Lastname:  googleUser.FamilyName,
				Email:     googleUser.Email,
			}
			u, err = delivery.userUsecase.CreateUser(req.Context(), &NewUserCred)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		token, err := delivery.userUsecase.CreateSessionJWT(c.Request.Context(), u, delivery.secretKey)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		//url := "http://localhost:3000" // TODO
		//
		//redirectURL := fmt.Sprintf(
		//	"%s?token=%s&refresh=%s",
		//	url,
		//	token.AccessToken,
		//	token.RefreshToken,
		//)
		//
		//http.Redirect(w, req, redirectURL , http.StatusFound)

		c.JSON(http.StatusOK, token)

	}
	return fn
}

func failure(c *gin.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.JSON(http.StatusNotFound, gin.H{"error": "google callback failure"})
	}
}

func (delivery *Delivery) LogoutUser(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery LogoutUser()")
	c.SetCookie("token", "", -1, "/", "http://localhost:3000", false, true) //TODO change to webapp url
	c.JSON(http.StatusOK, gin.H{"you have been successfully logged out": nil})

}

// LoginUserYandex -
func (delivery *Delivery) LoginUserYandex(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery LoginUserYandex()")
	oauth2ConfigYandex := &oauth2.Config{
		ClientID:     "c0ba17e8f61d47fdb0a978131d1c2d48",
		ClientSecret: "2e41418551724d918919ac09d4f6a1eb",
		Endpoint:     yandex.Endpoint,
		RedirectURL:  "http://localhost:8000/user/callbackYandex",
		Scopes:       []string{"email"},
	}
	c.Redirect(http.StatusTemporaryRedirect, oauth2ConfigYandex.AuthCodeURL("random")) //

	delivery.logger.Debug(yandex.Endpoint.TokenURL)
	c.JSON(http.StatusOK, yandex.Endpoint)
}

// CallbackYandex -
func (delivery *Delivery) CallbackYandex(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery CallbackYandex()")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
	//c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	//c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	//stateConfig := gologin.DefaultCookieConfig
	//gologinoauth2.StateHandler(stateConfig, )
	//c.JSON(http.StatusOK, gin.H{})

}
