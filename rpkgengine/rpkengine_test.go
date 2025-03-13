package rpkgengine_test

import (
	"testing"

	re "github.com/rsdate/rpkgengine/rpkgengine"
	rec "github.com/rsdate/rpkgengine/rpkgengineconfig"
	e "github.com/rsdate/utils/errors"
	"github.com/spf13/viper"
)

var (
	// f is the struct that holds the data from rpkg.build.yaml file
	f re.RpkgBuildFile = re.RpkgBuildFile{}
	// viper_instance is the viper instance (used to remove redundancy)
	viper_instance *viper.Viper = viper.GetViper()
	// errChecker is the error checker (with test mode on)
	errChecker = e.ErrChecker{
		ErrPrefix: "test error",
		PanicMode: "true",
		EM:        e.EM["eMre"],
		TestMode:  true,
	}
	Em = errChecker.EM
)

// This function is mainly for standalone testing (not recommended for production).
// It initializes the configuration and builds the package.
// It uses the rpkg-test package.
func TestRpkgengineBuild(t *testing.T) {
	var _, err = errChecker.CheckErr(Em[""], func() (any, error) {
		var _, err = rec.InitConfig("../../rpkg-test/")
		f.Authors = viper_instance.Get("authors").([]any)
		f.BuildCommands = viper_instance.Get("build_commands").([]any)
		f.BuildDeps = viper_instance.Get("build_deps").([]any)
		f.BuildWith = viper_instance.Get("build_with").(string)
		f.Deps = viper_instance.Get("deps").([]any)
		f.Name = viper_instance.Get("name").(string)
		f.Revision = viper_instance.Get("revision").(int)
		f.Version = viper_instance.Get("version").(string)
		err = re.Build("../rpkg-test", f, false)
		return nil, err
	})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
