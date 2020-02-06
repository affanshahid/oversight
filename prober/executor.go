package prober

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/affanshahid/oversight/prober/probe"
	"github.com/spf13/viper"
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
		if err != nil {
			fmt.Println(err)
			// TODO: what to do in case of error?
			continue
		}

		err = e.saveInRedis(data)
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(time.Duration(e.probe.GetConfig().Interval) * time.Millisecond)
	}
}

func (e *executor) saveInRedis(data *probe.RawSignal) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	redisClient := e.probe.GetRedisClient()
	return redisClient.LPush(viper.GetString("redis.to_process_queue"), bytes).Err()
}

// NewExecutor creates and returns a new Executor
func newExecutor(p probe.Probe) *executor {
	return &executor{probe: p}
}
