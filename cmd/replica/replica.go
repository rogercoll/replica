package main

import (
	"context"
	"flag"
	"strings"

	"github.com/rogercoll/replica/config"
	"github.com/rogercoll/replica/controller"
	_ "github.com/rogercoll/replica/plugins/backup/all"
	_ "github.com/rogercoll/replica/plugins/distributors/all"
	"github.com/sirupsen/logrus"
)

var fSampleConfig = flag.Bool("sample-config", false,
	"print out full sample configuration")

var fDistFilters = flag.String("dist-filter", "",
	"filter the distributors to enable, separator is :")

var fBckFilters = flag.String("bck-filter", "",
	"filter the backups to enable, separator is :")

var fConfig = flag.String("config", "",
	"configuration file path")

const (
	version = "v0.0.1"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func runReplica(distFilters, backupFilters []string) {
	log.WithFields(logrus.Fields{"version": version}).Info("Starting Replica!")

	// If no other options are specified, load the config file and run.
	c := config.NewConfig(log.WithFields(logrus.Fields{"version": version}))
	err := c.LoadConfig(*fConfig)
	if err != nil {
		log.Fatal(err)
	}
	ctr, err := controller.NewController(c)
	if err != nil {
		log.Fatal(err)
	}
	err = ctr.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	distFilters, bckFilters := []string{}, []string{}

	if *fDistFilters != "" {
		distFilters = strings.Split(strings.TrimSpace(*fDistFilters), ":")
	}
	if *fBckFilters != "" {
		bckFilters = strings.Split(strings.TrimSpace(*fBckFilters), ":")
	}
	switch {
	case *fSampleConfig:
		config.PrintSampleConfig(bckFilters, distFilters)
		return
	}
	run(
		distFilters,
		bckFilters,
	)
}
