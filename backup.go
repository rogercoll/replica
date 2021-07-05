package replica

type Backup interface {
	PluginDescriber
	//Should we return error or the backup should have its own logger
	//Returns a list with the paths of the backups
	Do() ([]string, error)
}
