package rpkgengineconfig

// Total lines in this file: 21
import (
	"github.com/spf13/viper"
)

// InitConfig reads in the rpkg.build.yaml file from projectPath.
func InitConfig(projectPath string) (int, error) {
	viper.SetConfigName("rpkg.build")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(projectPath)
	if err := viper.ReadInConfig(); err != nil {
		return 1, err
	}
	return 0, nil
}
