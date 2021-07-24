package replica

type Distributor interface {
	PluginDescriber
	Save([]Backup) (int64, error)
}
