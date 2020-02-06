package prober

import (
	"fmt"
	"github.com/affanshahid/oversight/prober/probe"
	"github.com/affanshahid/oversight/prober/registrar"

	"github.com/jinzhu/gorm"
)

// Prober controls all probing logic
type Prober struct {
	db        *gorm.DB
	executors []*executor
}

// Start starts the prober system
func (p *Prober) Start() error {
	var configs []*probe.Config
	p.db.Find(&configs)

	for _, config := range configs {
		probe, err := registrar.NewProbe(config)
		if err != nil {
			fmt.Printf("unable to start probe using config id %s skipping", config.ID)
			continue
		}

		executor := newExecutor(probe)
		p.executors = append(p.executors, executor)
		executor.Start()
	}

	return nil
}

// Stop stops the prober system
func (p *Prober) Stop() {
	for _, executor := range p.executors {
		executor.Stop()
	}
}

// New creates new prober
func New(db *gorm.DB) *Prober {
	if db == nil {
		panic("DB can not be nil")
	}

	return &Prober{db: db}
}
