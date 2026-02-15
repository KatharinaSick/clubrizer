package app

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Url string `yaml:"url"`
	} `yaml:"database"`

	OAuth struct {
		AccessToken struct {
			Issuer     string `yaml:"issuer"`
			HeaderName string `yaml:"headerName"`
			SecretKey  string `yaml:"secretKey"`
		} `yaml:"accessToken"`
		RefreshToken struct {
			Issuer     string `yaml:"issuer"`
			CookieName string `yaml:"cookieName"`
			SecretKey  string `yaml:"secretKey"`
			Secure     bool   `yaml:"secure"`
			SameSite   string `yaml:"sameSite"`
		} `yaml:"refreshToken"`
		User struct {
			Key string `yaml:"key"`
		} `yaml:"user"`
		Google struct {
			ClientId     string `yaml:"clientId"`
			ClientSecret string `yaml:"clientSecret"`
		} `yaml:"google"`
	} `yaml:"oauth"`

	Cors struct {
		AllowedOrigins []string `yaml:"allowedOrigins"`
		AllowedHeaders []string `yaml:"allowedHeaders"`
	} `yaml:"cors"`
}

func ParseConfig() (*Config, error) {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		return nil, err
	}
	replaced := os.ExpandEnv(string(data))
	cfg := &Config{}
	err = yaml.Unmarshal([]byte(replaced), cfg)
	return cfg, err
}
