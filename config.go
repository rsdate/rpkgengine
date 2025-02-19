package main

import (
	"github.com/spf13/viper"
)

// Description: initConfig initializes the configuration
//
// Parameters: It takes the project path
//
// Returns: It returns an integer and an error. The integer is the exit code of the configuration initialization process (1 or 0) and the error is any error that occurred during the configuration initialization process
func initConfig(projectPath string) (int, error) {
	viper.SetConfigName("rpkg.build")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(projectPath)
	if err := viper.ReadInConfig(); err != nil {
		return 1, err
	}
	return 0, nil
}
