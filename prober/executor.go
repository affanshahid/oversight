package prober

import (
	"fmt"
	"github.com/affanshahid/oversight/prober/probe"
	"time"
)

type executor struct {
	probe       probe.Probe
	stopChannel chan bool
}

// Start the probe
func (e *executor) Start() error {
	err := e.probe.BeforeStart()
	if err != nil {
		return err
	}

	e.stopChannel = make(chan bool)
	go e.process()

	err = e.probe.AfterStart()
	if err != nil {
		e.Stop()
		return err
	}

	return nil
}

// Stop stops the probe
func (e *executor) Stop() {
	e.stopChannel <- true
	e.probe.BeforeShutdown()
	close(e.stopChannel)
	e.stopChannel = nil
}

func (e *executor) shouldContinueProcessing() bool {
	select {
	case <-e.stopChannel:
		return false
	default:
		return true
	}
}

func (e *executor) process() {
	for {
		shouldContinue := e.shouldContinueProcessing()

		if !shouldContinue {
			return
		}

		data, err := e.probe.Fetch()
		if err == nil {
			fmt.Printf("[%s] Got data length: %d\n", e.probe.GetConfig().ID, len(data))
		}
		// TODO: what to do in case of error?

		time.Sleep(time.Duration(e.probe.GetConfig().Interval) * time.Millisecond)
	}
}

// NewExecutor creates and returns a new Executor
func newExecutor(p probe.Probe) *executor {
	return &executor{probe: p}
}
