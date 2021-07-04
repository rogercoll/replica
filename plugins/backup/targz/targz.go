package targz

import (
	"fmt"
	"os"

	"github.com/rogercoll/replica"
	"github.com/rogercoll/replica/plugins/backup"
)

type TarGz struct {
	Paths      []string `toml:"paths"`
	TimeFormat string   `toml:"timeformat"`
}

const sampleConfig = `
  ## Create Tar Gz files for given paths (files or directories)
  # This is expected to be the list of absolute host path
  paths = []
  # Time format to be used in the backup file names (https://golang.org/pkg/time/#Time)
  # Default is layoutISO = "2006-01-02"
  timeformat = "2006-01-02T15:04:05Z07:00" 
`

// SampleConfig returns a sample TOML section to illustrate configuration
// options.
func (t *TarGz) SampleConfig() string {
	return sampleConfig
}

func (t *TarGz) Description() string {
	return "Generate Tar Gz files for given paths"
}

func (t *TarGz) Do() ([]*os.File, error) {
	fmt.Println("Doing backup")
	return nil, nil
}

func init() {
	backup.Add("targz", func() replica.Backup {
		return &TarGz{TimeFormat: "2006-01-02"}
	})
}
