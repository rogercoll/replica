package replica

import "io"

//This is a more generic backup instance. It allows for any type of backups, like docker volumes.
//For example, a targz backup would return its filename in the Name() function and the bytes of the targz in the Data() one.
type Backup interface {
	Name() string
	Data() io.Reader
	Error() error
}

//More generic, not tied to a persistance file on disk and add clean functionality
type BackupSystem interface {
	PluginDescriber
	Do(<-chan Backup)
	Clean()
}
