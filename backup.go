package replica

import "io"

type Backup interface {
	PluginDescriber
	//Should we return error or the backup should have its own logger
	//Returns a list with the paths of the backups
	Do() ([]string, error)
}

//This is a more generic backup instance. It allows for any type of backups, like docker volumes.
//For example, a targz backup would return its filename in the Name() function and the bytes of the targz in the Data() one.
type Backup2 interface {
	Name() string
	Data() io.Reader
}

//More generic, not tied to a persistance file on disk and add clean functionality
type BackupSystem interface {
	PluginDescriber
	Do()
	GetBackups(<-chan Backup2)
	Clean()
}
