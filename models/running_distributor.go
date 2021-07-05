package models

import (
	"github.com/rogercoll/replica"
	"github.com/sirupsen/logrus"
)

type RunningDistributor struct {
	Dist replica.Distributor
	Name string
	Log  *logrus.Entry
}

func NewRunningDistributor(name string, dist replica.Distributor) *RunningDistributor {
	logger := logrus.New().WithFields(logrus.Fields{"distributor": name})
	return &RunningDistributor{
		Dist: dist,
		Log:  logger,
		Name: name,
	}
}
