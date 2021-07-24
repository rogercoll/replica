package models

import (
	"github.com/rogercoll/replica"
	"github.com/sirupsen/logrus"
)

type RunningBackup struct {
	Backup replica.BackupSystem
	Name   string
	log    *logrus.Entry
}

func NewRunningBackup(name string, backup replica.BackupSystem, _log *logrus.Entry) *RunningBackup {
	_log.WithFields(logrus.Fields{"backup": name})
	return &RunningBackup{
		Backup: backup,
		Name:   name,
		log:    _log,
	}
}
