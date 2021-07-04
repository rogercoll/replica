package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/rogercoll/replica/config"
	"github.com/rogercoll/replica/controller"
	_ "github.com/rogercoll/replica/plugins/auth/all"
	_ "github.com/rogercoll/replica/plugins/backup/all"
)

var fSampleConfig = flag.Bool("sample-config", false,
	"print out full sample configuration")

var fAuthFilters = flag.String("auth-filter", "",
	"filter the authorizators to enable, separator is :")

var fBckFilters = flag.String("bck-filter", "",
	"filter the backups to enable, separator is :")

var fConfig = flag.String("config", "",
	"configuration file path")

func runReplica(authFilters, backupFilters []string) {
	version := "v0.0.1"
	log.Printf("I! Starting Replica %s", version)

	// If no other options are specified, load the config file and run.
	c := config.NewConfig()
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
	authFilters, bckFilters := []string{}, []string{}

	if *fAuthFilters != "" {
		authFilters = strings.Split(strings.TrimSpace(*fAuthFilters), ":")
	}
	if *fBckFilters != "" {
		bckFilters = strings.Split(strings.TrimSpace(*fBckFilters), ":")
	}
	switch {
	case *fSampleConfig:
		config.PrintSampleConfig(bckFilters, authFilters)
		return
	}
	run(
		authFilters,
		bckFilters,
	)
}
