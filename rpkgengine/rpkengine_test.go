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
	errCheckerTest = e.ErrChecker{
		ErrPrefix: "test error",
		PanicMode: "true",
		EM:        e.EM["eMre"],
		TestMode:  true,
	}
	Em = errCheckerTest.EM
)

// This function is mainly for standalone testing (not recommended for production).
// It initializes the configuration and builds the package.
// It uses the rpkg-test package.
func TestRpkgengineBuild(t *testing.T) {
	var _, err = errCheckerTest.CheckErr(Em[""], func() (any, error) {
		var _, err = rec.InitConfig("../../rpkg-test/")
		f = re.InitVars(viper_instance)
		err = re.Build("../../rpkg-test", f, false, re.ErrCheckerBuild)
		return nil, err
	})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
