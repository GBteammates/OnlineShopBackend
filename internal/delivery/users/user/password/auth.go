package password

import (
	"OnlineShopBackend/internal/models"
	"crypto/sha1"
	"fmt"
	"log"
	"unicode"

	"github.com/caarlos0/env/v6"
)

type SaltSHA struct {
	Salt string `json:"salt" env:"SALT" envDefault:"sjdhkashdsw823rgfeg"`
}

type Credentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewPasswordConfig() (*SaltSHA, error) {
	//var configPathHash = "./internal/delivery/user/password/hash.json"

	var cfg = SaltSHA{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("can't load environment variables: %s", err)
	}

	//data, err := os.ReadFile(configPathHash)
	//if err != nil {
	//	log.Fatalf("cannot read the file: %s", err)
	//}
	//
	//if err = json.Unmarshal(data, &cfg); err != nil {
	//	log.Fatalf("cannot unmarshal: %s", err)
	//}

	return &cfg, nil
}

func ValidationCheck(user models.User) error {
	if user.Email == "" && user.Firstname == "" && user.Lastname == "" {
		return fmt.Errorf("empty filed")
	}
	if len(user.Password) < 5 {
		return fmt.Errorf("password is too short")
	}
	for _, char := range user.Password {
		if !unicode.IsDigit(char) && !unicode.Is(unicode.Latin, char) {
			return fmt.Errorf("password should contain lathin letter or numbers only")
		}
	}
	return nil
}

func GeneratePasswordHash(password string) string {
	cfg, err := NewPasswordConfig()
	if err != nil {
		return ""
	}

	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(cfg.Salt)))
}

func SanitizePassword(user *models.User) {
	user.Password = ""
}
