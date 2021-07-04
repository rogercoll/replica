package replica

import "os"

type Authenticator interface {
	PluginDescriber
	Save([]*os.File) (int64, error)
}
