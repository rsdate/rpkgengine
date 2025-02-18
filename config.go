package main

import (
	"github.com/spf13/viper"
)

func initConfig() (int, error) {
	viper.SetConfigName("rpkg.build")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return 1, err
	}
	return 0, nil
}
