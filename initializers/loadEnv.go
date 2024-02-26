package initializers

import (
	"fmt"

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

	GOOGLE_MAPS_API_KEY               string `mapstructure:"GOOGLE_MAPS_API_KEY"`
	STRIPE_SECRET_KEY                 string `mapstructure:"STRIPE_SECRET_KEY"`
	STRIPE_WEBHOOK_SIGNING_SECRET_KEY string `mapstructure:"STRIPE_WEBHOOK_SIGNING_SECRET_KEY"`
	//APP CONFIG
	APP_URL    string `mapstructure:"APP_URL"`
	APP_SECRET string `mapstructure:"APP_SECRET"`

	CLIENT_ORDER_SERVICE_FEE    uint64 `mapstructure:"ORDER_SERVICE_FEE"`
	RESPONDER_ORDER_SERVICE_FEE uint64 `mapstructure:"RESPONDER_ORDER_SERVICE_FEE"`

	RESPONDENT_ORDER_CHARGE_AMOUNT uint64 `mapstructure:"RESPONDENT_ORDER_CHARGE_AMOUNT"`

	SERVER_ADDRESS string `mapstructure:"SERVER_ADDRESS"`
	SERVER_PORT    string `mapstructure:"SERVER_PORT"`
	PUBLIC_PREFIX  string `mapstructure:"PUBLIC_PREFIX"`
	// UPLOAD_DIR   string `mapstructure:"PUBLIC_PREFIX"`
	//

	ASSET_DIR              string `mapstructure:"ASSET_DIR"`
	UPLOAD_DIR             string `mapstructure:"UPLOAD_DIR"`
	USER_AVATAR_UPLOAD_DIR string `mapstructure:"USER_AVATAR_UPLOAD_DIR"`

	STRIPE_RESPONDENT_PURCHASE_PRODUCT_ID     string `mapstructure:"STRIPE_RESPONDENT_PURCHASE_PRODUCT_ID"`
	STRIPE_RESPONDENT_SUBSCRIPTION_PRODUCT_ID string `mapstructure:"STRIPE_RESPONDENT_SUBSCRIPTION_PRODUCT_ID"`
	//

	STRIPE_RESPONDENT_PURCHASE_PRODUCT_SUCCESS_CALLBACK_URL string `mapstructure:"STRIPE_RESPONDENT_PURCHASE_PRODUCT_SUCCESS_CALLBACK_URL"`
	STRIPE_RESPONDENT_PURCHASE_PRODUCT_FAILURE_CALLBACK_URL string `mapstructure:"STRIPE_RESPONDENT_PURCHASE_PRODUCT_FAILURE_CALLBACK_URL"`

	//

	SENDGRID_API_KEY    string `mapstructure:"SENDGRID_API_KEY"`
	SENDGRID_FROM_EMAIL string `mapstructure:"SENDGRID_FROM_EMAIL"`
	SENDGRID_FROM_NAME  string `mapstructure:"SENDGRID_FROM_NAME"`
}

var loaded bool

func LoadConfig(paths ...string) (config *Config, err error) {

	if loaded {
		return CONFIG, nil
	}
	var path string
	if len(paths) > 0 {
		path = paths[0]
	} else {
		path = "."
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
	CONFIG = config

	stripe.Key = config.STRIPE_SECRET_KEY

	loaded = true

	return config, err
}

func (c *Config) GetServerAddress() string {
	// return c.SERVER_ADDRESS + ":" + c.SERVER_PORT
	serverAddr := c.SERVER_ADDRESS
	if c.SERVER_PORT != "" {
		serverAddr += ":" + c.SERVER_PORT
	}
	return serverAddr
}

func (c *Config) ResolvePublicPath(filePath string) string {
	if c.PUBLIC_PREFIX != "" {
		filePath = fmt.Sprintf("%s/%s", c.PUBLIC_PREFIX, filePath)
	}
	return filePath
}
