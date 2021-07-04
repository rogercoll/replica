package config

import (
	"fmt"
	"strings"

	"github.com/rogercoll/replica"
	"github.com/rogercoll/replica/plugins/auth"
)

var authHeader = `
###############################################################################
#                            AUTH PLUGINS                                    #
###############################################################################
`

// PrintSampleConfig prints the sample config
func PrintSampleConfig(authFilters []string) {
	if len(authFilters) > 0 {
		fmt.Printf(authHeader)
		for _, filter := range authFilters {
			creator := auth.Auths[filter]
			auth := creator()
			printConfig(filter, auth, "auths", false)
		}
	} else {
		fmt.Println("No Auth filters found")
	}
}

func printConfig(name string, p replica.PluginDescriber, op string, commented bool) {
	comment := ""
	if commented {
		comment = "# "
	}
	fmt.Printf("\n%s# %s\n%s[[%s.%s]]", comment, p.Description(), comment,
		op, name)

	config := p.SampleConfig()
	if config == "" {
		fmt.Printf("\n%s  # no configuration\n\n", comment)
	} else {
		lines := strings.Split(config, "\n")
		for i, line := range lines {
			if i == 0 || i == len(lines)-1 {
				fmt.Print("\n")
				continue
			}
			fmt.Print(strings.TrimRight(comment+line, " ") + "\n")
		}
	}
}
