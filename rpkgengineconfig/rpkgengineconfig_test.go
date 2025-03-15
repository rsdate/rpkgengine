package rpkgengineconfig_test

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
		BuildCommands: []any{"python3.13 -m build"},
	}
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

func TestInitConfig(t *testing.T) {
	var f, err = errCheckerTest.CheckErr(Em[""], func() (any, error) {
		var _, err = rec.InitConfig("../../rpkg-test/")
		f = re.InitVars(viper_instance)
		return f, err
	})
	F := f.(re.RpkgBuildFile)
	if err != nil {
		t.Errorf("Error: %v", err)
	} else if !reflect.DeepEqual(F, exampleRpkgBuildFile) {
		t.Errorf("Error: %s", "f != exampleRpkgBuildFile")
	}
}
