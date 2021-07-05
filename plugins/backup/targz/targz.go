package targz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

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

func (t *TarGz) Do() ([]string, error) {
	var backups []string
	for _, aPath := range t.Paths {
		body, err := ioutil.ReadFile(aPath)
		if err != nil {
			log.Fatalln(err)
		}
		outputName := path.Base(aPath) + "-" + time.Now().Format(t.TimeFormat)
		file, err := os.Create("/tmp/" + outputName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		writer, err := gzip.NewWriterLevel(file, gzip.BestCompression)
		if err != nil {
			log.Fatalln(err)
		}
		defer writer.Close()

		tw := tar.NewWriter(writer)
		defer tw.Close()
		if body != nil {
			hdr := &tar.Header{
				Name: outputName,
				Mode: int64(0644),
				Size: int64(len(body)),
			}
			if err := tw.WriteHeader(hdr); err != nil {
				fmt.Println(err)
			}
			if _, err := tw.Write(body); err != nil {
				fmt.Println(err)
			}
		}
		backups = append(backups, file.Name())
	}
	return backups, nil
}

func init() {
	backup.Add("targz", func() replica.Backup {
		return &TarGz{TimeFormat: "2006-01-02"}
	})
}
