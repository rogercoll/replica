package replica

import "os"

type Backup interface {
	PluginDescriber
	//Should we return error or the backup should have its own logger
	Do() ([]*os.File, error)
}
