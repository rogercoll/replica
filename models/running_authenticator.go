package models

import (
	"github.com/rogercoll/replica"
	"github.com/sirupsen/logrus"
)

type RunningAuthenticator struct {
	Auth replica.Authenticator
	Name string
	log  *logrus.Entry
}

func NewRunningAuthenticator(name string, auth replica.Authenticator) *RunningAuthenticator {
	logger := logrus.New().WithFields(logrus.Fields{"auth": name})
	return &RunningAuthenticator{
		Auth: auth,
		log:  logger,
	}
}
