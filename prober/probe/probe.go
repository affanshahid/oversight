package probe

import (
	"fmt"
)

// Probe represents all functionality required in a custom probe
type Probe interface {
	BeforeStart() error     // Called everytime before the probe is started
	AfterStart() error      // Called everytime after a successful start
	Fetch() (string, error) // Called to fetch data
	BeforeShutdown() error  // Called before the probe is shutdown
	GetConfig() *Config     // Returns the config
}

// BaseProbe encapsulates default probe functionality
type BaseProbe struct {
	Config *Config
}

// BeforeStart does nothing
func (*BaseProbe) BeforeStart() error {
	fmt.Println("Base BeforeStart()")
	return nil
}

// AfterStart does nothing
func (*BaseProbe) AfterStart() error {
	fmt.Println("Base AfterStart()")
	return nil
}

// BeforeShutdown does nothing
func (*BaseProbe) BeforeShutdown() error {
	fmt.Println("Base BeforeShutdown()")
	return nil
}

// GetConfig returns the config
func (b *BaseProbe) GetConfig() *Config {
	return b.Config
}
