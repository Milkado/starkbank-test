package helpers

import (
	"fmt"

	"github.com/spf13/viper"
)

func Env(key string) string {
	viper.SetConfigFile("../.env")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("error getting the provided key from .env: ", err.Error())
		panic("viper error")
	}

	return viper.GetString(key)
}