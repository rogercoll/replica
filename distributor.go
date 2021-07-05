package replica

type Distributor interface {
	PluginDescriber
	Save([]string) (int64, error)
}
