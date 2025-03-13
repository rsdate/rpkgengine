package rpkgengine

import e "github.com/rsdate/utils/errors"

var (
	errChecker = e.ErrChecker{
		ErrPrefix: "build error",
		PanicMode: "true",
		EM:        e.EM["eMre"],
		TestMode:  false,
	}
	Em = errChecker.EM
)
