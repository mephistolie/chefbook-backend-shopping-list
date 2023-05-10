package config

import (
	"github.com/mephistolie/chefbook-backend-common/log"
)

const (
	EnvDebug = "debug"
	EnvProd  = "production"
)

type Config struct {
	Environment *string
	Port        *int
	LogsPath    *string
}

func (c Config) Validate() error {
	if *c.Environment != EnvProd {
		*c.Environment = EnvDebug
	}
	return nil
}

func (c Config) Print() {
	log.Infof("SERVICE CONFIGURATION\n"+
		"Environment: %v\n"+
		"Port: %v\n"+
		"Logs path: %v\n\n"+
		*c.Environment, *c.Port, *c.LogsPath,
	)
}
