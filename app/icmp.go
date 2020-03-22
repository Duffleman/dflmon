package app

import (
	"strings"
	"time"

	"dflmon/config"

	sdk "github.com/andygrunwald/cachet"
	log "github.com/sirupsen/logrus"
	"github.com/sparrc/go-ping"
)

// PacketsToSend is the number of packets to send
const PacketsToSend = 2

// Timeout is how long to wait before saying the ICMP failed
const Timeout = 15 * time.Second

// Interval is the wait time between the packets to send
const Interval = 5 * time.Millisecond

func (a *App) doICMP(job *config.Job) int {
	pinger, err := ping.NewPinger(job.Host)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			log.Warnf("no such host, configuration error for host %s", job.Host)
		}

		return sdk.ComponentStatusMajorOutage
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
		return sdk.ComponentStatusMajorOutage
	// some packets returned
	case stats.PacketsRecv < stats.PacketsSent:
		l.Infof("packet loss on host %s", job.Host)
		return sdk.ComponentStatusPartialOutage
	// all packets retured
	case stats.PacketsRecv == stats.PacketsSent:
		l.Infof("successfully pinged host %s", job.Host)
		return sdk.ComponentStatusOperational
	}

	l.Warnf("unknown state found")

	return sdk.ComponentStatusUnknown
}
