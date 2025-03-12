package main

// Total lines in the package: 448
// Total lines in this file: 30
import (
	re "github.com/rsdate/rpkgengine/rpkgengine"
	"github.com/spf13/viper"
)

var (
	// f is the struct that holds the data from rpkg.build.yaml file
	f re.RpkgBuildFile = re.RpkgBuildFile{}
	// viper_instance is the viper instance (used to remove redundancy)
	viper_instance *viper.Viper = viper.GetViper()
)

// This function is mainly for standalone testing (not recommended for production).
// It initializes the configuration and builds the package.
// It uses the rpkg-test package.
func main() {
	initConfig("../rpkg-test/")
	f.Authors = viper_instance.Get("authors").([]any)
	f.BuildCommands = viper_instance.Get("build_commands").([]any)
	f.BuildDeps = viper_instance.Get("build_deps").([]any)
	f.BuildWith = viper_instance.Get("build_with").(string)
	f.Deps = viper_instance.Get("deps").([]any)
	f.Name = viper_instance.Get("name").(string)
	f.Revision = viper_instance.Get("revision").(int)
	f.Version = viper_instance.Get("version").(string)
	re.Build("../rpkg-test", f, false)
}
