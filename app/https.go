package app

import (
	"fmt"
	"strings"

	"dflmon"
	"dflmon/config"

	log "github.com/sirupsen/logrus"
)

func (a *App) doHTTPS(job *config.Job) error {
	url := fmt.Sprintf("https://%s", job.Host)

	res, err := a.Client.Get(url)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			log.Warnf("no such host, configuration error for host %s", job.Host)
		}

		log.Infof("cannot connect to host %s", job.Host)

		return dflmon.ErrMajorOutage
	}

	if res.StatusCode > 400 {
		log.Infof("cannot connect to host %s", job.Host)
		return dflmon.ErrMajorOutage
	}

	log.Infof("successfully connected to host %s", job.Host)

	return nil
}
