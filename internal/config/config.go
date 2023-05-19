package config

import (
	"github.com/mephistolie/chefbook-backend-common/log"
	"time"
)

const (
	EnvDev  = "develop"
	EnvProd = "production"
)

type Config struct {
	Environment *string
	Port        *int
	LogsPath    *string

	ShoppingList ShoppingList

	Firebase Firebase
	Database Database
	Amqp     Amqp
}

type ShoppingList struct {
	MaxShoppingListsCount *int
	KeyTtl                *time.Duration
	CheckSubscription     *bool
}

type Firebase struct {
	Credentials *string
}

type Database struct {
	Host     *string
	Port     *int
	User     *string
	Password *string
	DBName   *string
}

type Amqp struct {
	Host     *string
	Port     *int
	User     *string
	Password *string
	VHost    *string
}

func (c Config) Validate() error {
	if *c.Environment != EnvProd {
		*c.Environment = EnvDev
	}
	return nil
}

func (c Config) Print() {
	log.Infof("SHOPPING LIST SERVICE CONFIGURATION\n"+
		"Environment: %v\n"+
		"Port: %v\n"+
		"Logs path: %v\n\n"+
		"Max shopping lists count: %v\n"+
		"Check subscription: %v\n\n"+
		"Database host: %v\n"+
		"Database port: %v\n"+
		"Database name: %v\n\n"+
		"MQ host: %v\n"+
		"MQ port: %v\n"+
		"MQ vhost: %v\n\n",
		*c.Environment, *c.Port, *c.LogsPath,
		*c.ShoppingList.MaxShoppingListsCount, *c.ShoppingList.CheckSubscription,
		*c.Database.Host, *c.Database.Port, *c.Database.DBName,
		*c.Amqp.Host, *c.Amqp.Port, *c.Amqp.VHost,
	)
}
