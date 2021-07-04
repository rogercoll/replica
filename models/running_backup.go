package models

import (
	"github.com/rogercoll/replica"
	"github.com/sirupsen/logrus"
)

type RunningBackup struct {
	Backup replica.Backup
	Name   string
	log    *logrus.Entry
}

func NewRunningBackup(name string, backup replica.Backup) *RunningBackup {
	logger := logrus.New().WithFields(logrus.Fields{"backup": name})
	return &RunningBackup{
		Backup: backup,
		Name:   name,
		log:    logger,
	}
}
