package config

import (
	"github.com/spf13/viper"
)

var GlobalConfig Config

type Config struct {
	APISecret string `mapstructure:"API_SECRET"`
}

func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&GlobalConfig)
	return
}
