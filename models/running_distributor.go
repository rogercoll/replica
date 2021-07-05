package models

import (
	"github.com/rogercoll/replica"
	"github.com/sirupsen/logrus"
)

type RunningDistributor struct {
	Dist replica.Distributor
	Name string
	log  *logrus.Entry
}

func NewRunningDistributor(name string, dist replica.Distributor) *RunningDistributor {
	logger := logrus.New().WithFields(logrus.Fields{"distributor": name})
	return &RunningDistributor{
		Dist: dist,
		log:  logger,
	}
}
