package ssh

import (
	"io"
	"log"
	"os"
	"path"
	"sync"

	"github.com/rogercoll/replica"
	"github.com/rogercoll/replica/plugins/distributors"
)

var (
	defaultDestination = os.ExpandEnv("${HOME}/backups/")
)

type Local struct {
	Destinations []string `toml:"destinations"`
}

const sampleConfig = `
  ## Copy backup files to other directories
  # Absolute direcotries path to store the backup files
  # Default is $HOME/backups
  destination = [] 
`

// SampleConfig returns a sample TOML section to illustrate configuration
// options.
func (l *Local) SampleConfig() string {
	return sampleConfig
}

func (l *Local) Description() string {
	return "Store backupfiles to a local directories"
}

func (l *Local) Save(files []string) (int64, error) {
	var total int64
	for _, file := range files {
		var wg sync.WaitGroup
		for _, localDir := range l.Destinations {
			wg.Add(1)
			go func(ldir string) {
				defer wg.Done()
				newFile, err := os.Create(ldir + path.Base(file))
				if err != nil {
					log.Fatal(err)
				}
				defer newFile.Close()
				oldFile, err := os.Open(file)
				if err != nil {
					log.Fatal(err)
				}
				defer oldFile.Close()
				bytesCopied, err := io.Copy(newFile, oldFile)
				if err != nil {
					log.Fatal(err)
				}
				total += bytesCopied
			}(localDir)
		}
		wg.Wait()
	}
	return total, nil
}

func init() {
	distributors.Add("local", func() replica.Distributor {
		return &Local{Destinations: []string{defaultDestination}}
	})
}
