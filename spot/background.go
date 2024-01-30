package spot

import (
	"time"
)

type backgroundRunner interface {
	start()
	stop()
}

type invalidationRunner struct {
	quit    chan struct{}
	garages *garageManager
}

func (r *invalidationRunner) start() {
	r.quit = make(chan struct{})
	go r.run()
}

func (r *invalidationRunner) stop() {
	r.quit <- struct{}{}
}

func (r *invalidationRunner) run() {
	for {
		select {
		case <-time.After(30 * time.Minute):
			r.garages.invalidateOldUpdates()
		case <-r.quit:
			return
		}
	}
}
