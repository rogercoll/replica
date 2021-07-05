package replica

type Authenticator interface {
	PluginDescriber
	Save([]string) (int64, error)
}
