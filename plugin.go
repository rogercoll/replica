package replica

type PluginDescriber interface {
	// SampleConfig returns the default configuration of the Processor
	SampleConfig() string

	// Description returns a one-sentence description on the Processor
	Description() string
}
