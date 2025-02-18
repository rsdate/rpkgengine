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
	initConfig("../rpkg-test/")
	f.Authors = viper_instance.Get("authors").([]interface{})
	f.BuildCommands = viper_instance.Get("build_commands").([]interface{})
	f.BuildDeps = viper_instance.Get("build_deps").([]interface{})
	f.BuildWith = viper_instance.Get("build_with").(string)
	f.Deps = viper_instance.Get("deps").([]interface{})
	f.Name = viper_instance.Get("name").(string)
	f.Revision = viper_instance.Get("revision").(int)
	f.Version = viper_instance.Get("version").(string)
	re.Build("../rpkg-test", f, true)
}
