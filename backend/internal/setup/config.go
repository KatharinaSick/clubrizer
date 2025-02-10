package setup

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
	}

	OAuth struct {
		SecretKey []byte `yaml:"secretKey"`
		Issuer    string `yaml:"issuer"`
		Google    struct {
			ClientId     string `yaml:"clientId"`
			ClientSecret string `yaml:"clientSecret"`
		} `yaml:"google"`
	} `yaml:"oauth"`
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
