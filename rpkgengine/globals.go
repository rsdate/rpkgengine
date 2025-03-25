package rpkgengine

import (
	e "github.com/rsdate/utils/errors"
)

var (
	conf            string
	mirror          string = "RPKG_MIRROR"
	Emre                   = e.EM["eMre"]
	ErrCheckerBuild        = e.ErrChecker{
		ErrPrefix: "build error",
		PanicMode: "true",
		EM:        Emre,
		TestMode:  false,
	}
)
