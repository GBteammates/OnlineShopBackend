/*
 * Backend for Online Shop
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var regexpEmail = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

type User struct {
	ID        uuid.UUID   `json:"id"`
	Firstname string      `json:"firstname,omitempty"`
	Lastname  string      `json:"lastname,omitempty"`
	Password  string      `json:"password,omitempty"`
	Email     string      `json:"email,omitempty"`
	Address   UserAddress `json:"address,omitempty"`
	Rights    Rights      `json:"rights"`
}

func (user *User) ValidationCheck(logger *zap.Logger) error {
	logger.Debug("Enter in models user ValidationCheck()")
	if user.Email == "" && user.Firstname == "" && user.Lastname == "" {
		return fmt.Errorf("empty fields")
	}
	if len(user.Firstname) > 100 || len(user.Lastname) > 100 {
		return fmt.Errorf("user name or user lastname too long")
	}
	if !regexpEmail.MatchString(strings.ToLower(user.Email)) {
		return errors.New("invalid email format")
	}
	if len(user.Password) < 5 {
		return fmt.Errorf("password is too short")
	}
	if len(user.Password) > 16 {
		return fmt.Errorf("password is too long")
	}
	for _, char := range user.Password {
		if !unicode.IsDigit(char) && !unicode.Is(unicode.Latin, char) {
			return fmt.Errorf("password should contain latin letter or numbers only")
		}
	}
	logger.Info("Validation success")
	return nil
}

func (user *User) GeneratePasswordHash(logger *zap.Logger) (string, error) {
	logger.Debug("Enter in models.User GeneratePasswordHash()")
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return "", err
	}
	logger.Info("Generation Password hash success")
	return string(bytes), nil
}

func (user *User) CheckPasswordHash(password string, logger *zap.Logger) bool {
	logger.Sugar().Debugf("Enter in models user CheckPasswordHash() with args: password: %s, logger", password)
	fmt.Printf("user password is: %v", user.Password)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	fmt.Printf("err compare hash and password is %v", err)
	return err == nil
}
