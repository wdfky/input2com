package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Debug    bool   `mapstructure:"debug"`
	Baudrate int    `mapstructure:"baudrate"`
	TtyPath  string `mapstructure:"ttyPath"`
	Server   struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	MouseConfigDict map[string]map[byte]string `mapstructure:"mouseConfigDict"`
}

var Cfg *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Cfg)
	if err != nil {
		panic(err)
	}
}
