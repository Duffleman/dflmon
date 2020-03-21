package app

import (
	"fmt"
	"strings"

	"dflmon"
	"dflmon/config"

	log "github.com/sirupsen/logrus"
)

func (a *App) doHTTPS(job *config.Job, validate bool) error {
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

		return dflmon.ErrMajorOutage
	}

	if res.StatusCode >= 400 {
		log.WithFields(log.Fields{
			"statusCode": res.StatusCode,
			"status":     res.Status,
		}).Infof("cannot connect to host %s", job.Host)
		return dflmon.ErrMajorOutage
	}

	log.Infof("successfully connected to host %s", job.Host)

	return nil
}
