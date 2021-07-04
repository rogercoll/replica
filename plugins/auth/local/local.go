package ssh

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/rogercoll/replica"
	"github.com/rogercoll/replica/plugins/auth"
)

var (
	defaultDestination = os.ExpandEnv("${HOME}/backups/")
)

type Local struct {
	Destination string `toml:"destination"`
}

const sampleConfig = `
  ## Store files in a local path
  # Absolute direcotry path to store the backup files
  # Default is $HOME/backups
  destination = "" 
`

// SampleConfig returns a sample TOML section to illustrate configuration
// options.
func (l *Local) SampleConfig() string {
	return sampleConfig
}

func (l *Local) Description() string {
	return "Store backupfiles to a local directory"
}

func (l *Local) Save(files []*os.File) (int64, error) {
	var total int64
	for _, file := range files {
		newFile, err := os.Create(l.Destination + path.Base(file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		defer newFile.Close()

		bytesCopied, err := io.Copy(newFile, file)
		if err != nil {
			log.Fatal(err)
		}
		total += bytesCopied
	}
	return total, nil
}

func init() {
	auth.Add("local", func() replica.Authenticator {
		return &Local{Destination: defaultDestination}
	})
}
