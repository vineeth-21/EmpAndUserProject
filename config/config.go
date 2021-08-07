package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() {
	log.Println("loadconfig called.....")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("config not found......")
		log.Panic("config not found......")
	}
}
