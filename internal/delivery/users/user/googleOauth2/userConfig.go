package googleOauth2

import (
	"encoding/json"
	"log"
	"os"
)

type UserConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_uri"`
}

func NewUserConfig() (*UserConfig, error) {
	var configPathGoogle = "./static/config/clientSecretGoogle.json"

	var cfg = UserConfig{}

	data, err := os.ReadFile(configPathGoogle)
	if err != nil {
		log.Fatalf("cannot read the file: %s", err)
	}

	if err = json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("cannot unmarshal: %s", err)
	}

	return &cfg, nil
}
