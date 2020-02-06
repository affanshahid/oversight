package probe

// RawSignal is what is pushed to redis by probes
type RawSignal struct {
	ProbeConfigID string
	Data          interface{}
}
