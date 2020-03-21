package app

import (
	"dflmon"
	"dflmon/config"

	sdk "github.com/andygrunwald/cachet"
)

func (a *App) HandleSuccess(job config.Job) error {
	if job.ComponentID == 0 {
		return dflmon.ErrNoComponent
	}

	_, _, err := a.Cachet.Components.Update(job.ComponentID, &sdk.Component{
		Status: sdk.ComponentStatusOperational,
	})
	if err != nil {
		return err
	}

	return nil
}
