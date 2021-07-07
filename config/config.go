package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/naoina/toml"
	"github.com/naoina/toml/ast"
	"github.com/rogercoll/replica"
	"github.com/rogercoll/replica/models"
	"github.com/rogercoll/replica/plugins/backup"
	"github.com/rogercoll/replica/plugins/distributors"
	"github.com/sirupsen/logrus"
)

type Config struct {
	toml          *toml.Config
	log           *logrus.Entry
	DistFilters   []string
	BackupFilters []string
	Distributors  []*models.RunningDistributor
	Backups       []*models.RunningBackup
}

func NewConfig(_log *logrus.Entry) *Config {
	return &Config{
		DistFilters:   make([]string, 0),
		BackupFilters: make([]string, 0),
		log:           _log,
		Distributors:  make([]*models.RunningDistributor, 0),
		Backups:       make([]*models.RunningBackup, 0),
	}
}

// Try to find a default config file at these locations (in order):
//   1. $REPLICA_CONFIG_PATH
//   2. $HOME/.replica/replica.conf
//   3. /etc/replica/replica.conf
//
func getDefaultConfigPath() (string, error) {
	envfile := os.Getenv("REPLICA_CONFIG_PATH")
	homefile := os.ExpandEnv("${HOME}/.replica/replica.conf")
	etcfile := "/etc/replica/replica.conf"
	for _, path := range []string{envfile, homefile, etcfile} {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	// if we got here, we didn't find a file in a default location
	return "", fmt.Errorf("No config file specified, and could not find one"+
		" in $REPLICA_CONFIG_PATH, %s, or %s", homefile, etcfile)
}

// LoadConfig loads the given config file and applies it to c
func (c *Config) LoadConfig(path string) error {
	var err error
	if path == "" {
		if path, err = getDefaultConfigPath(); err != nil {
			return err
		}
		c.log.WithFields(logrus.Fields{"path": path}).Info("I! Using config file")
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Error loading config file %s: %w", path, err)
	}

	if err = c.LoadConfigData(data); err != nil {
		return fmt.Errorf("Error loading config file %s: %w", path, err)
	}
	return nil
}

// LoadConfigData loads TOML-formatted config data
func (c *Config) LoadConfigData(data []byte) error {
	tbl, err := toml.Parse(data)
	if err != nil {
		return fmt.Errorf("Error parsing data: %s", err)
	}
	// Parse all the rest of the plugins:
	for name, val := range tbl.Fields {
		subTable, ok := val.(*ast.Table)
		if !ok {
			return fmt.Errorf("invalid configuration, error parsing field %q as table", name)
		}

		switch name {
		case "backup":
			for pluginName, pluginVal := range subTable.Fields {
				switch pluginSubTable := pluginVal.(type) {
				case *ast.Table:
					if err = c.addBackup(pluginName, pluginSubTable); err != nil {
						return fmt.Errorf("error parsing %s, %w", pluginName, err)
					}
				case []*ast.Table:
					for _, t := range pluginSubTable {
						if err = c.addBackup(pluginName, t); err != nil {
							return fmt.Errorf("error parsing %s, %w", pluginName, err)
						}
					}
				default:
					return fmt.Errorf("Unsupported config format: %s",
						pluginName)
				}
			}
		case "distributor":
			for pluginName, pluginVal := range subTable.Fields {
				switch pluginSubTable := pluginVal.(type) {
				// legacy [inputs.cpu] support
				case *ast.Table:
					if err = c.addDistributor(pluginName, pluginSubTable); err != nil {
						return fmt.Errorf("error parsing %s, %w", pluginName, err)
					}
				case []*ast.Table:
					for _, t := range pluginSubTable {
						if err = c.addDistributor(pluginName, t); err != nil {
							return fmt.Errorf("error parsing %s, %w", pluginName, err)
						}
					}
				default:
					return fmt.Errorf("Unsupported config format: %s",
						pluginName)
				}
			}
		}
	}
	return nil
}

func (c *Config) addBackup(name string, table *ast.Table) error {
	if len(c.BackupFilters) > 0 && !sliceContains(name, c.BackupFilters) {
		return nil
	}
	creator, ok := backup.Backups[name]
	if !ok {
		return fmt.Errorf("Undefined but requested input: %s", name)
	}
	bck := creator()
	if err := c.toml.UnmarshalTable(table, bck); err != nil {
		return err
	}

	rp := models.NewRunningBackup(name, bck, c.log)
	c.Backups = append(c.Backups, rp)
	return nil
}

func (c *Config) addDistributor(name string, table *ast.Table) error {
	if len(c.DistFilters) > 0 && !sliceContains(name, c.DistFilters) {
		return nil
	}
	creator, ok := distributors.Distributors[name]
	if !ok {
		return fmt.Errorf("Undefined but requested distributor: %s", name)
	}
	dist := creator()
	if err := c.toml.UnmarshalTable(table, dist); err != nil {
		return err
	}

	rp := models.NewRunningDistributor(name, dist, c.log)
	c.Distributors = append(c.Distributors, rp)
	return nil
}

var distHeader = `
###############################################################################
#                            DISTRIBUTOR PLUGINS                              #
###############################################################################
`

var bckHeader = `
###############################################################################
#                           BACKUP PLUGINS                                    #
###############################################################################
`

// PrintSampleConfig prints the sample config
func PrintSampleConfig(bckFilters, distFilters []string) {
	if len(bckFilters) > 0 {
		fmt.Printf(bckHeader)
		for _, filter := range bckFilters {
			creator := backup.Backups[filter]
			backup := creator()
			printConfig(filter, backup, "backup", false)
		}
	} else {
		fmt.Println("No Backup filters found")
	}
	if len(distFilters) > 0 {
		fmt.Printf(distHeader)
		for _, filter := range distFilters {
			creator := distributors.Distributors[filter]
			dist := creator()
			printConfig(filter, dist, "distributor", false)
		}
	} else {
		fmt.Println("No Distributor filters found")
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

func sliceContains(name string, list []string) bool {
	for _, b := range list {
		if b == name {
			return true
		}
	}
	return false
}
