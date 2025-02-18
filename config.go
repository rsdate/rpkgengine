package main

import (
	"github.com/spf13/viper"
)

func initConfig(projectPath string) (int, error) {
	viper.SetConfigName("rpkg.build")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(projectPath)
	if err := viper.ReadInConfig(); err != nil {
		return 1, err
	}
	return 0, nil
}
