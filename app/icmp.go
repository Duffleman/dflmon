package app

import (
	"strings"

	"dflmon"
	"dflmon/config"

	log "github.com/sirupsen/logrus"
	"github.com/sparrc/go-ping"
)

const PacketsToSend = 2
const TimeoutInSeconds = 5

func (a *App) doICMP(job *config.Job) error {
	pinger, err := ping.NewPinger(job.Host)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			log.Warnf("no such host, configuration error for host %s", job.Host)
		}

		return dflmon.ErrMajorOutage
	}

	pinger.SetPrivileged(true)

	pinger.Count = PacketsToSend
	pinger.Timeout = TimeoutInSeconds

	pinger.Run()

	stats := pinger.Statistics()

	switch {
	// no packets returned
	case stats.PacketsRecv == 0:
		log.Infof("cannot ping host %s", job.Host)
		return dflmon.ErrMajorOutage
	// some packets returned
	case stats.PacketsRecv < stats.PacketsSent:
		log.Infof("packet loss on host %s", job.Host)
		return dflmon.ErrPartialOutage
	// all packets retured
	case stats.PacketsRecv == stats.PacketsSent:
		log.Infof("successfully pinged host %s", job.Host)
		return nil
	}

	log.WithFields(log.Fields{
		"stats": stats,
	}).Warnf("unknown state found")

	return dflmon.ErrUnknownState
}
