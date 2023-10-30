package helloworldmetricsprocessor

import (
	"fmt"
)

// Config represents the receiver config settings within the collector's config.yaml
type Config struct {
	ExampleParameterAttr string `mapstructure:"attribute"`
}

// Validate checks if the receiver configuration is valid
func (cfg *Config) Validate() error {
	if cfg.ExampleParameterAttr == "" {
		return fmt.Errorf("attribute parameter not set")
	}
	return nil
}
