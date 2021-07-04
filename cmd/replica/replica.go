package main

import (
	"flag"
	"strings"

	"github.com/rogercoll/replica/config"
	_ "github.com/rogercoll/replica/plugins/auth/all"
)

var fSampleConfig = flag.Bool("sample-config", false,
	"print out full sample configuration")

var fAuthFilters = flag.String("auth-filter", "",
	"filter the authorizators to enable, separator is :")

func main() {
	flag.Parse()
	authFilters := []string{}
	if *fAuthFilters != "" {
		authFilters = strings.Split(strings.TrimSpace(*fAuthFilters), ":")
	}
	switch {
	case *fSampleConfig:
		config.PrintSampleConfig(authFilters)
		return
	}
}
