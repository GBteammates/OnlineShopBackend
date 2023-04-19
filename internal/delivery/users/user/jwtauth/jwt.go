package jwtauth

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type Payload struct {
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	UserId uuid.UUID `json:"userId"`
	jwt.StandardClaims
}

func NewJWT(payload Payload) (string, error) {
	key, err := NewJWTKeyConfig()
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)
	return token.SignedString([]byte(key.Key))
}

func NewRefreshToken() (string, error) {
	refreshToken := make([]byte, 32)
	t := rand.NewSource(time.Now().Unix())
	r := rand.New(t)

	_, err := r.Read(refreshToken)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", refreshToken), nil
}

func CreateSessionJWT(ctx context.Context, user *models.User) (Token, error) {
	payload := Payload{
		Email:  user.Email,
		Role:   user.Rights.Name,
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	accessToken, err := NewJWT(payload)
	if err != nil {
		return Token{}, fmt.Errorf("unable to create a token")
	}

	refreshToken, err := NewRefreshToken()
	if err != nil {
		return Token{}, fmt.Errorf("unable to create a refresh token")
	}

	token := Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	//if err = usecase.userStore.SaveSession(ctx, refreshToken, payload.ExpiresAt); err != nil {
	//	return Token{}, fmt.Errorf("unable to save session")
	//}

	return token, nil
}
