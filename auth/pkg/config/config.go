package config

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Port          string `mapstructure:"PORT"`
	SecretKey     string `mapstructure:"SECRET_KEY"`
	AuthSvcUrl    string `mapstructure:"AUTH_SVC_URL"`
	EmailSvcUrl   string `mapstructure:"EMAIL_SVC_URL"`
	ProductSvcUrl string `mapstructure:"PRODUCT_SVC_URL"`
	OrderSvcUrl   string `mapstructure:"ORDER_SVC_URL"`
	DBUrl         string `mapstructure:"DB_URL"`
	HOST          string `mapstructure:"HOST"`
	FromEmail     string `mapstructure:"EMAIL"`
	FromPassword  string `mapstructure:"PASSWORD"`
	RedisAddr     string `mapstructure:"REDIS_ADDRESS"`
}

func LoadConfig() (c Config, err error) {
	viper.AddConfigPath("./pkg/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err = viper.Unmarshal(&c)

	log.Info().Msg("Load Config....")
	return
}
