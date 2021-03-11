package services

import "github.com/tkanos/gonfig"

// region ======== TYPES =================================================================

// conf unexported configuration data holder struct
type conf struct {

	// Database configuration
	Host     string
	User     string
	Pass     string
	Database string
}

// SvcConfig unexported configuration service struct
type SvcConfig struct {
	Path string `string:"Path to the config YAML file"`
	conf `conf:"Configuration object"`
}

// endregion =============================================================================

// NewConfig create a new configuration service.
//
// - path [string] ~ Path to the configuration .yaml file
func NewConfig(path string) *SvcConfig {
	c := conf{}

	err := gonfig.GetConf(path, &c) // getting the conf
	if err != nil {                 // error check
		panic(err)
	}

	return &SvcConfig{path, c}
}
