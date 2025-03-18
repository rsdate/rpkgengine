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
		err = re.Build("../../rpkg-test", f, false, errCheckerTest)
		return nil, err
	})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestDownloadPackage(t *testing.T) {
	var _, err = errCheckerTest.CheckErr(Em[""], func() (any, error) {
		err := re.DownloadPackage("./libdaemon-0.0.1.tar.gz", "https://rsdate.github.io/projects/libdaemon-0.0.1.tar.gz", errCheckerTest)
		return nil, err
	})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestBuildPackage(t *testing.T) {
	var _, err = errCheckerTest.CheckErr(Em[""], func() (any, error) {
		var _, err = rec.InitConfig("../../rpkg-test/")
		err = re.BuildPackage("../../rpkg-test/", errCheckerTest)
		return nil, err
	})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

/* Not just yet
func TestInstallPackage(t *testing.T) {
	var _, err = errCheckerTest.CheckErr(Em[""], func() (any, error) {
		err := re.InstallPackage("./libdaemon-0.0.1.tar.gz", "lidaemon-0.0.1.tar.gz", "libdaemon-0.0.1", "false", true, errCheckerTest)
		return nil, err
	})
	if err == nil {
		t.Errorf("Error: no error returned")
	}
}
*/
