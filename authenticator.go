package replica

import "io"

type Authenticator interface {
	PluginDescriber
	Save(io.Writer) error
}
