package registrar

import (
	"errors"
	"fmt"
	"github.com/affanshahid/oversight/prober/probe"
)

// ProbeFactory is a type alias for a function which
// creates a Probe given a config
type ProbeFactory = func(conf *probe.Config) (probe.Probe, error)

var registry = make(map[string]ProbeFactory)

var (
	// ErrNotRegistered is returned by NewProbe when no matching key is found
	ErrNotRegistered = errors.New("key not registered")
)

// Register registers a probe factory against a descriminator string
func Register(key string, factory ProbeFactory) {
	_, exists := registry[key]

	if exists {
		panic("attempting to re-register probe factory with key '" + key + "'")
	}

	if factory == nil {
		panic("attempting to register a nil probe factory with key '" + key + "'")
	}

	registry[key] = factory
}

// NewProbe creates a new probe using config.Descriminator as the key
// to retrieve the factory from the registry
func NewProbe(config *probe.Config) (probe.Probe, error) {
	factory, exists := registry[config.Descriminator]

	if !exists {
		return nil, fmt.Errorf("key '%s' not registered: %w", config.Descriminator, ErrNotRegistered)
	}

	return factory(config)
}
