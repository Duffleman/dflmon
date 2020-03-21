package app

import (
	"errors"
	"net"

	"dflmon"

	sdk "github.com/andygrunwald/cachet"
)

var ErrDNS = &net.DNSError{}

func (a *App) HandleError(err JobWrap) error {
	if err.Job.ComponentID == 0 {
		return dflmon.ErrNoComponent
	}

	var switchTo int

	switch err.Err {
	case dflmon.ErrMajorOutage:
		switchTo = sdk.ComponentStatusMajorOutage
	case dflmon.ErrPartialOutage:
		switchTo = sdk.ComponentStatusPartialOutage
	case dflmon.ErrPerformanceIssue:
		switchTo = sdk.ComponentStatusPerformanceIssues
	default:
		return errors.New("unknown error receiver")
	}

	_, _, newErr := a.Cachet.Components.Update(err.Job.ComponentID, &sdk.Component{
		Status: switchTo,
	})
	if newErr != nil {
		return newErr
	}

	return nil
}
