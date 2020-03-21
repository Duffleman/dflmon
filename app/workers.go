package app

import (
	"math/rand"
	"sync"
	"time"

	"dflmon/config"

	log "github.com/sirupsen/logrus"
)

var (
	MaxRand = 10
	MinRand = 2
)

type JobWrap struct {
	Job config.Job
	Err error
}

func (a *App) StartWorkers() error {
	wg := &sync.WaitGroup{}

	errCh := make(chan JobWrap)

	go a.MessageHandlerWorker(errCh)

	for i, job := range a.Config.Jobs {
		log.Infof("starting worker %d/%d", i+1, len(a.Config.Jobs))
		wg.Add(1)
		go a.StartWorker(wg, errCh, job)
		wait := time.Duration(rand.Intn((MaxRand - MinRand) + MinRand))
		log.Infof("waiting for %d seconds", wait)
		time.Sleep(wait * time.Second)
	}

	wg.Wait()

	return nil
}

func (a *App) StartWorker(wg *sync.WaitGroup, errCh chan JobWrap, job *config.Job) {
	defer wg.Done()

	for {
		switch job.Type {
		case "icmp":
			err := a.doICMP(job)
			errCh <- JobWrap{*job, err}
		case "https":
			err := a.doHTTPS(job)
			errCh <- JobWrap{*job, err}
		default:
			log.Warnf("job type not implemented %s", job.Type)
			return
		}

		time.Sleep(job.Interval * time.Second)
	}
}
