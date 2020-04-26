package app

import (
	"fmt"
	"strings"

	"dflmon/config"

	sdk "github.com/andygrunwald/cachet"
	log "github.com/sirupsen/logrus"
)

var allowedCodes = map[int]struct{}{
	200: struct{}{},
	204: struct{}{},
	302: struct{}{},
	401: struct{}{},
}

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

	if _, ok := allowedCodes[res.StatusCode]; !ok {
		l.Infof("cannot connect to host %s", job.Host)

		return sdk.ComponentStatusMajorOutage
	}

	l.Infof("successfully connected to host %s", job.Host)

	return sdk.ComponentStatusOperational
}
