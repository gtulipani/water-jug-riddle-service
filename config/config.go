package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config represents main config.
type Config struct {
	HTTPPort string
}

// InitConfig: loads required configuration
func InitConfig() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	c := Config{
		HTTPPort: v.GetString(httpPort),
	}

	if err := validateConfig(v); err != nil {
		return nil, err
	}

	return &c, nil
}

func validateConfig(viper *viper.Viper) error {
	mandatoryVariables := []string{
		httpPort,
	}

	for _, v := range mandatoryVariables {
		if viper.Get(v) == nil {
			return fmt.Errorf("missing mandatory environment variable: %s", v)
		}
	}

	return nil
}