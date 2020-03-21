package app

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (a *App) HandleMessage(err JobWrap) {
	var newErr error

	if err.Err == nil {
		newErr = a.HandleSuccess(err.Job)
		if newErr != nil {
			log.Warn(fmt.Errorf("cannot update component %s: %w", err.Job.Name, newErr))
		}
		return
	}

	newErr = a.HandleError(err)
	if newErr != nil {
		log.Warn(fmt.Errorf("cannot update component %s: %w", err.Job.Name, newErr))
	}
}

func (a *App) MessageHandlerWorker(ch chan JobWrap) {
	for {
		select {
		case err, ok := <-ch:
			if !ok {
				return
			}

			a.HandleMessage(err)
		}
	}
}
