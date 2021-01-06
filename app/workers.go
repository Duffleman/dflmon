package app

import (
	"math/rand"
	"sync"
	"time"

	"dflmon/config"

	log "github.com/sirupsen/logrus"
)

var (
	// MaxRand is the maximum time to wait before starting a job
	MaxRand = 10
	// MinRand is the minimum time to wait before starting a job
	MinRand = 2
)

type jobWrap struct {
	Job     config.Job
	Outcome int
}

// StartWorkers starts all workers, 1 per job
func (a *App) StartWorkers() error {
	wg := &sync.WaitGroup{}

	jobCh := make(chan jobWrap)

	go a.messageHandlerWorker(jobCh)

	for i, job := range a.Config.Jobs {
		log.Infof("starting worker %d/%d", i+1, len(a.Config.Jobs))
		wg.Add(1)
		go a.startWorker(wg, jobCh, job)
		wait := time.Duration(rand.Intn((MaxRand - MinRand) + MinRand))
		log.Infof("waiting for %d seconds", wait)
		time.Sleep(wait * time.Second)
	}

	wg.Wait()

	return nil
}

func (a *App) startWorker(wg *sync.WaitGroup, jobCh chan jobWrap, job *config.Job) {
	defer wg.Done()

	var outcome int

	for {
		switch job.Type {
		case "icmp":
			outcome = a.doICMP(job)
		case "https":
			outcome = a.doHTTPS(job, true)
		case "https-novalidate":
			outcome = a.doHTTPS(job, false)
		case "http":
			outcome = a.doHTTP(job, true)
		case "tcp":
			outcome = a.doTCP(job)
		default:
			log.Warnf("job type not implemented %s", job.Type)
			return
		}

		jobCh <- jobWrap{*job, outcome}

		time.Sleep(job.Interval * time.Second)
	}
}
