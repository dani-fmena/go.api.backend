package utils

import "github.com/tkanos/gonfig"

// region ======== TYPES =================================================================

// conf unexported configuration schema holder struct
type conf struct {

	// Database configuration
	Addr	 string
	Host     string
	Port     string
	User     string
	Pass     string
	Database string

	// Migrations directory
	MigrationDir string

	// Environment
	Debug bool
}

// SvcConfig unexported configuration service struct
type SvcConfig struct {
	Path string `string:"Path to the config YAML file"`
	conf `conf:"Configuration object"`
}
// endregion =============================================================================

// NewSvcConfig create a new configuration service.
//
// - path [string] ~ Path to the configuration .yaml file
func NewSvcConfig(path string) *SvcConfig {
	c := conf{}

	err := gonfig.GetConf(path, &c) 			// getting the conf
	if err != nil {                 // error check
		panic(err)
	}

	return &SvcConfig{path, c} // We are using struct composition here. Hence the anonymous field (https://golangbot.com/inheritance/)
}
