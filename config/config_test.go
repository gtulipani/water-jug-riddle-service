package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	routeTests := []struct {
		name                 string
		environmentVariables map[string]string
		expectedError        error
		output               *Config
	}{
		{
			name: "error without httpPort",
			environmentVariables: map[string]string{},
			expectedError: fmt.Errorf("missing mandatory environment variable: %s", httpPort),
		},
	}

	for _, tt := range routeTests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Unsetenv(httpPort)

			for k, v := range tt.environmentVariables {
				_ = os.Setenv(k, v)
			}

			c, err := InitConfig()
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.output, c)
		})
	}
}