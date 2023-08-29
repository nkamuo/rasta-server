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
}

func LoadConfig(path string) (config Config, err error) {
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

	return config, err
}
