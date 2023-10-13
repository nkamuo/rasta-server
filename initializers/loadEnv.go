package initializers

import (
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v74"
)

var CONFIG *Config

type Config struct {
	DBHost         string `mapstructure:"MYSQL_HOST"`
	DBUserName     string `mapstructure:"MYSQL_USER"`
	DBUserPassword string `mapstructure:"MYSQL_PASSWORD"`
	DBName         string `mapstructure:"MYSQL_DATABASE"`
	DBPort         string `mapstructure:"MYSQL_PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`

	GOOGLE_MAPS_API_KEY string `mapstructure:"GOOGLE_MAPS_API_KEY"`
	STRIPE_SECRET_KEY   string `mapstructure:"STRIPE_SECRET_KEY"`
	//APP CONFIG
	APP_URL    string `mapstructure:"APP_URL"`
	APP_SECRET string `mapstructure:"APP_SECRET"`

	CLIENT_ORDER_SERVICE_FEE    uint64 `mapstructure:"ORDER_SERVICE_FEE"`
	RESPONDER_ORDER_SERVICE_FEE uint64 `mapstructure:"RESPONDER_ORDER_SERVICE_FEE"`

	SERVER_ADDRESS string `mapstructure:"SERVER_ADDRESS"`
	SERVER_PORT    string `mapstructure:"SERVER_PORT"`
}

var loaded bool

func LoadConfig(path string) (config Config, err error) {

	if loaded {
		return *CONFIG, nil
	}
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	CONFIG = &config

	stripe.Key = config.STRIPE_SECRET_KEY

	loaded = true

	return config, err
}
