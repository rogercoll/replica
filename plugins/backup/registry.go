package backup

import "github.com/rogercoll/replica"

type Creator func() replica.Backup

var Backups = map[string]Creator{}

func Add(name string, creator Creator) {
	Backups[name] = creator
}
