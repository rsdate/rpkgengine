package rpkgengine_test

import (
	"reflect"
	"testing"

	re "github.com/rsdate/rpkgengine/rpkgengine"
	rec "github.com/rsdate/rpkgengine/rpkgengineconfig"
	e "github.com/rsdate/utils/errors"
	"github.com/spf13/viper"
)

var (
	// f is the struct that holds the data from rpkg.build.yaml file
	f re.RpkgBuildFile = re.RpkgBuildFile{}
	// exampleRpkgBuildFile is the example rpkg.build.yaml file
	exampleRpkgBuildFile = re.RpkgBuildFile{
		Name:          "test_pkg",
		Version:       "0.0.1",
		Revision:      0,
		Authors:       []any{"Rohan Date <rohan.s.date@icloud.com> AUTHOR", "Jane Doe <jane.doe@test.com> MAINTAINER"},
		Deps:          []any{"none"},
		BuildDeps:     []any{"build@latest"},
		BuildWith:     "python3.13",
		BuildCommands: []any{"python3 -m build"},
	}
	// viper_instance is the viper instance (used to remove redundancy)
	viper_instance *viper.Viper = viper.GetViper()
	// errChecker is the error checker (with test mode on)
	errCheckerTest = e.ErrChecker{
		ErrPrefix: "test error",
		PanicMode: "true",
		EM:        e.EM["eMIer"],
		TestMode:  true,
	}
	Em = errCheckerTest.EM
)

func TestInstallDeps(t *testing.T) {
	deps := []any{"build@latest"}
	installCommands := map[string][]string{
		"latest":  {"install", "--upgrade"},
		"version": {"install"},
	}
	var _, err = errCheckerTest.CheckErr("tsterr1", func() (any, error) {
		err := re.InstallDeps(deps, true, "pip", installCommands, "==")
		return nil, err
	})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	Deps := []any{"build@1.2.2.post1"}
	var _, Err = errCheckerTest.CheckErr("tsterr1", func() (any, error) {
		err := re.InstallDeps(Deps, true, "pip", installCommands, "==")
		return nil, err
	})
	if err != nil {
		t.Errorf("Error: %v", Err)
	}

}

func TestInitVars(t *testing.T) {
	var _, _ = rec.InitConfig("../rpkg-test")
	f := re.InitVars(viper_instance)
	if !reflect.DeepEqual(f, exampleRpkgBuildFile) {
		t.Errorf("Error: f != exampleRpkgBuildFile")
	}
}

func TestRpkgengineBuild(t *testing.T) {
	var _, err = errCheckerTest.CheckErr("tsterr1", func() (any, error) {
		var _, err = rec.InitConfig("../rpkg-test/")
		f = re.InitVars(viper_instance)
		err = re.Build("../rpkg-test", f, false, errCheckerTest)
		return nil, err
	})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestDownloadPackage(t *testing.T) {
	var _, err = errCheckerTest.CheckErr("tsterr1", func() (any, error) {
		err := re.DownloadPackage("./libdaemon-0.0.1.tar.gz", "https://rsdate.github.io/projects/libdaemon-0.0.1.tar.gz", errCheckerTest)
		return nil, err
	})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestBuildPackage(t *testing.T) {
	var _, err = errCheckerTest.CheckErr("tsterr1", func() (any, error) {
		var _, err = rec.InitConfig("../rpkg-test/")
		err = re.BuildPackage("../rpkg-test/", errCheckerTest, viper_instance)
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
