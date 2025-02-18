package main

import (
	re "github.com/rsdate/rpkgengine/rpkgengine"
	"github.com/spf13/viper"
)

var (
	f              re.RpkgBuildFile = re.RpkgBuildFile{}
	viper_instance                  = viper.GetViper()
)

func main() {
	initConfig()
	f.Authors = []string{viper_instance.Get("authors").(string)}
	f.BuildCommands = viper_instance.Get("buildCommands").([]string)
	f.BuildDeps = viper_instance.Get("buildDeps").([]string)
	f.BuildWith = viper_instance.Get("buildWith").(string)
	f.Deps = viper_instance.Get("deps").([]string)
	f.Name = viper_instance.Get("name").(string)
	f.Revision = viper_instance.Get("revision").(int)
	f.Version = viper_instance.Get("version").(string)
	re.Build("../rpkg-test", f)
}
