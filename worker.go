package signalutils

import (
	"time"

	"github.com/sirupsen/logrus"
)

//Worker utility for launching Go routines that will loop over a function
//the max frequency of calls to this function is limited and
//the actual frequency is measured
type Worker struct {
	minFreq         float64
	maxFreq         float64
	ticker          *time.Ticker
	done            chan (bool)
	step            StepFunc
	stopOnErr       bool
	name            string
	active          bool
	CurrentFreq     float64
	CurrentStepTime time.Duration
}

//StepFunc function interface for the application that will be
//called in a loop
type StepFunc func() error

//StartWorker launches a Go routine looping in this step function limiting by maxFreq
//if the function is being run in a frequency less than minFreq, a logrus.Info log will show this
//this situation happens when the function is too slow
func StartWorker(name string, step StepFunc, minFreq float64, maxFreq float64, stopOnErr bool) *Worker {
	c := &Worker{
		name:      name,
		minFreq:   minFreq,
		maxFreq:   maxFreq,
		done:      make(chan bool),
		ticker:    time.NewTicker(time.Duration((float64(time.Second) / maxFreq))),
		step:      step,
		stopOnErr: stopOnErr,
		active:    false,
	}
	logrus.Tracef("Starting goroutine for %s", name)
	go c.run()
	return c
}

//Stop stops Go routine loop
func (c *Worker) Stop() {
	c.done <- true
}

func (c *Worker) run() {
	c.active = true
	for {
		loopStart := time.Now()
		select {
		case <-c.done:
			c.active = false
			return
		case <-c.ticker.C:
			stepStart := time.Now()
			err := c.step()
			c.CurrentStepTime = time.Since(stepStart)
			loopElapsed := time.Since(loopStart)
			c.CurrentFreq = float64(1) / loopElapsed.Seconds()
			logrus.Debugf("%s: STEP time=%d ms; loop freq=%.2f", c.name, c.CurrentStepTime.Milliseconds(), c.CurrentFreq)
			if err != nil {
				logrus.Infof("%s: STEP err=%s", c.name, err)
				if c.stopOnErr {
					c.active = false
					return
				}
			}
			if c.CurrentFreq < c.minFreq {
				logrus.Infof("%s: STEP too slow; loop freq=%.2f (min=%.2f)", c.name, c.CurrentFreq, c.minFreq)
			}
		}
	}
}
