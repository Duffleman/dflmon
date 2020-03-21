package app

import (
	"net/http"

	"dflmon/cachet"
	"dflmon/config"
)

type App struct {
	Client *http.Client
	Config *config.Config
	Cachet *cachet.Client
}

func (a *App) SyncWithCachet() error {
	components, err := a.Cachet.ListAllComponents()
	if err != nil {
		return err
	}

OUT:
	for _, job := range a.Config.Jobs {
		for _, component := range components {
			if component.Name == job.ComponentName {
				job.ComponentID = component.ID
				continue OUT
			}
		}
	}

	return nil
}
