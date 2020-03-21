package dflmon

import (
	"errors"
)

var (
	ErrPerformanceIssue = errors.New("performance issues")
	ErrPartialOutage    = errors.New("partial outage")
	ErrMajorOutage      = errors.New("major outage")
	ErrUnknownState     = errors.New("unknown state")

	ErrNoComponent = errors.New("no component matched")
)
