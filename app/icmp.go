package app

import (
	"strings"
	"time"

	"dflmon"
	"dflmon/config"

	log "github.com/sirupsen/logrus"
	"github.com/sparrc/go-ping"
)

const PacketsToSend = 2
const Timeout = 15 * time.Second
const Interval = 5 * time.Millisecond

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
	pinger.Timeout = Timeout
	pinger.Interval = Interval

	pinger.Run()

	stats := pinger.Statistics()

	l := log.WithFields(log.Fields{
		"stats": stats,
	})

	switch {
	// no packets returned
	case stats.PacketsRecv == 0:
		l.Infof("cannot ping host %s", job.Host)
		return dflmon.ErrMajorOutage
	// some packets returned
	case stats.PacketsRecv < stats.PacketsSent:
		l.Infof("packet loss on host %s", job.Host)
		return dflmon.ErrPartialOutage
	// all packets retured
	case stats.PacketsRecv == stats.PacketsSent:
		l.Infof("successfully pinged host %s", job.Host)
		return nil
	}

	l.Warnf("unknown state found")

	return dflmon.ErrUnknownState
}
