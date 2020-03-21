package app

import (
	"strings"

	"dflmon"
	"dflmon/config"

	log "github.com/sirupsen/logrus"
	"github.com/sparrc/go-ping"
)

const PacketsToSend = 2

func (a *App) doICMP(job *config.Job) error {
	pinger, err := ping.NewPinger(job.Host)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			log.Warnf("no such host, configuration error for host %s", job.Host)
		}

		return dflmon.ErrMajorOutage
	}

	pinger.Count = PacketsToSend

	pinger.Run()

	stats := pinger.Statistics()

	switch {
	case stats.PacketsRecv == PacketsToSend:
		log.Infof("successfully pinged host %s", job.Host)
		return nil
	case stats.PacketLoss > 0 && stats.PacketLoss < PacketsToSend:
		log.Infof("packet loss on host %s", job.Host)
		return dflmon.ErrPartialOutage
	case stats.PacketLoss == PacketsToSend:
		log.Infof("cannot ping host %s", job.Host)
		return dflmon.ErrMajorOutage
	}

	log.Warnf("unknown state found")

	return dflmon.ErrPartialOutage
}
