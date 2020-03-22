package app

import (
	"fmt"
	"strings"

	"dflmon/config"

	sdk "github.com/andygrunwald/cachet"
	log "github.com/sirupsen/logrus"
)

func (a *App) doHTTPS(job *config.Job, validate bool) int {
	url := fmt.Sprintf("https://%s", job.Host)

	c := a.Client

	if validate == false {
		c = a.ClientNoValidate
	}

	res, err := c.Get(url)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			log.Warnf("no such host, configuration error for host %s", job.Host)
		}

		log.WithFields(log.Fields{
			"error": err,
		}).Infof("cannot connect to host %s", job.Host)

		return sdk.ComponentStatusMajorOutage
	}

	l := log.WithFields(log.Fields{
		"statusCode": res.StatusCode,
		"status":     res.Status,
	})

	if res.StatusCode >= 400 {
		l.Infof("cannot connect to host %s", job.Host)

		return sdk.ComponentStatusMajorOutage
	}

	l.Infof("successfully connected to host %s", job.Host)

	return sdk.ComponentStatusOperational
}
