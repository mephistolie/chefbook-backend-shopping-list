package main

import (
	"flag"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/app"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/config"
	"github.com/peterbourgon/ff/v3"
	"os"
	"time"
)

func main() {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	cfg := config.Config{
		Environment: fs.String("environment", "debug", "service environment"),
		Port:        fs.Int("port", 8080, "service port"),
		LogsPath:    fs.String("logs-path", "", "logs file path"),

		ShoppingList: config.ShoppingList{
			MaxShoppingListsCount: fs.Int("max-shopping-lists", 5, "max shopping lists per owner count"),
			KeyTtl:                fs.Duration("shopping-list-key-ttl", 24*time.Hour, "shopping list key time to life"),
			CheckSubscription:     fs.Bool("check-subscription", true, "enable free subscription limits"),
		},

		ProfileService: config.Service{
			Addr: fs.String("profile-addr", "", "profile service address"),
		},

		RecipeService: config.Service{
			Addr: fs.String("recipe-addr", "", "recipe service address"),
		},

		Firebase: config.Firebase{
			Credentials: fs.String("firebase-credentials", "", "Firebase credentials JSON; leave empty to disable"),
		},

		Database: config.Database{
			Host:     fs.String("db-host", "localhost", "database host"),
			Port:     fs.Int("db-port", 5432, "database port"),
			User:     fs.String("db-user", "", "database user name"),
			Password: fs.String("db-password", "", "database user password"),
			DBName:   fs.String("db-name", "", "service database name"),
		},

		Amqp: config.Amqp{
			Host:     fs.String("amqp-host", "", "message broker host; leave empty to disable"),
			Port:     fs.Int("amqp-port", 5672, "message broker port"),
			User:     fs.String("amqp-user", "guest", "message broker user name"),
			Password: fs.String("amqp-password", "guest", "message broker user password"),
			VHost:    fs.String("amqp-vhost", "", "message broker virtual host"),
		},
	}
	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVars()); err != nil {
		panic(err)
	}

	err := cfg.Validate()
	if err != nil {
		panic(err)
	}

	app.Run(&cfg)
}
