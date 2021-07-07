package targz

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/rogercoll/replica"
	"github.com/rogercoll/replica/plugins/backup"
	"golang.org/x/sync/errgroup"
)

var (
	suffix = ".tar.gz"
)

type TarGz struct {
	Paths      []string `toml:"paths"`
	TimeFormat string   `toml:"timeformat"`
}

const sampleConfig = `
  ## Create Tar Gz files for given directories
  # This is expected to be the list of absolute host path directories
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

func compress(src string, buf io.Writer) error {
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)

	//maindir := path.Base(src)
	unusedPrefix := len(src) - len(path.Base(src))
	// walk through every file in the folder
	filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		// generate tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)
		header.Name = filepath.ToSlash(file[unusedPrefix:])

		// write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// if not a dir, write file content
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}
	// produce gzip
	if err := zr.Close(); err != nil {
		return err
	}
	return nil
}
func (t *TarGz) Do() ([]string, error) {
	backups := make([]string, len(t.Paths))
	g := new(errgroup.Group)
	for i, aPath := range t.Paths {
		aPath := aPath
		i := i
		g.Go(func() error {
			var buf bytes.Buffer
			err := compress(aPath, &buf)
			if err != nil {
				fmt.Println(err)
				return err
			}
			outputName := path.Base(aPath) + "-" + time.Now().Format(t.TimeFormat) + suffix
			file, err := os.Create("/tmp/" + outputName)
			if err != nil {
				fmt.Println(err)
				return err
			}
			defer file.Close()
			if _, err := io.Copy(file, &buf); err != nil {
				panic(err)
			}
			backups[i] = file.Name()
			return nil
		})
	}
	err := g.Wait()
	return backups, err
}

func init() {
	backup.Add("targz", func() replica.Backup {
		return &TarGz{TimeFormat: "2006-01-02"}
	})
}
