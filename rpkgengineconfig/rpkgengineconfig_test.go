package rpkgengineconfig_test

import (
	"testing"

	rec "github.com/rsdate/rpkgengine/rpkgengineconfig"
	e "github.com/rsdate/utils/errors"
)

var (
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
	var _, err = errCheckerTest.CheckErr(Em[""], func() (any, error) {
		var _, err = rec.InitConfig("../rpkg-test/")
		return nil, err
	})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
