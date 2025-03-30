package rpkgengine

import (
	e "github.com/rsdate/utils/errors"
)

const (
	installCommand string = "install"
	upgradeOption  string = "--upgrade"
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
	installCommandArray map[string][]string = map[string][]string{
		"latest":  {"install", "--upgrade"},
		"version": {"install"},
	}
)
