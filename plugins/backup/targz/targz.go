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

func (t *TarGz) Do() ([]*os.File, error) {
	var backups []*os.File
	fmt.Println(t.Paths)
	for _, aPath := range t.Paths {
		fmt.Println("Making backup of a file")
		body, err := ioutil.ReadFile(aPath)
		if err != nil {
			log.Fatalln(err)
		}
		outputName := path.Base(aPath) + "-" + time.Now().Format(t.TimeFormat)
		file, err := ioutil.TempFile("/tmp", outputName)
		if err != nil {
			log.Fatal(err)
		}
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
				println(err)
			}
			if _, err := tw.Write(body); err != nil {
				println(err)
			}
		}
		backups = append(backups, file)
	}
	return backups, nil
}

func init() {
	backup.Add("targz", func() replica.Backup {
		return &TarGz{TimeFormat: "2006-01-02"}
	})
}
