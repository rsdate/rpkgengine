package rpkgengine

import (
	e "github.com/rsdate/utils/errors"
	"github.com/spf13/viper"
)

var (
	conf            string
	mirror          string       = "RPKG_MIRROR"
	viper_instance  *viper.Viper = viper.GetViper()
	ErrCheckerBuild              = e.ErrChecker{
		ErrPrefix: "build error",
		PanicMode: "true",
		EM:        e.EM["eMre"],
		TestMode:  false,
	}
	Emre = e.EM["eMre"]
)
