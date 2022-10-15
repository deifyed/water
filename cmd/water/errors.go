package water

import "errors"

var (
	errMissingArguments = errors.New("no arguments provided. See --help for usage")
	errTargetNotExists  = errors.New("target does not exist")
)
