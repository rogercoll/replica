package ssh

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/rogercoll/replica"
	"github.com/rogercoll/replica/plugins/distributors"
	"golang.org/x/crypto/ssh"
)

type SSH struct {
	Hosts     []string `toml:"hosts"`
	Passwords []string `toml:"passwords"`
}

const sampleConfig = `
  ## Replicate files via SSH
  # This is expected to be the list of ssh config made of hostIP, port and user to do the replication, eg: "192.168.1.90:22:bob"
  hosts = []
  # Password to use with the corresponding config, if empty will use the loaded ssh keys
  usernames = []
`

// SampleConfig returns a sample TOML section to illustrate configuration
// options.
func (s *SSH) SampleConfig() string {
	return sampleConfig
}

func (s *SSH) Description() string {
	return "Connect via SSH and copy files via SCP"
}

func (s *SSH) Save(files []string) (int64, error) {
	var totalBytes int64
	for i, c := range s.Hosts {
		config := strings.Split(c, ":")
		if len(config) != 3 {
			return 0, errors.New("Invalid host format")
		}
		clientConfig := &ssh.ClientConfig{
			User: config[2],
			Auth: []ssh.AuthMethod{
				ssh.Password(s.Passwords[i]),
			},
		}
		client, err := ssh.Dial("tcp", config[0]+":"+config[1], clientConfig)
		if err != nil {
			return 0, err
		}
		session, err := client.NewSession()
		if err != nil {
			return 0, err
		}
		defer session.Close()
		r, err := session.StdoutPipe()
		if err != nil {
			return 0, err
		}

		for _, fileName := range files {
			file, err := os.Open(fileName)
			if err != nil {
				return 0, err
			}
			defer file.Close()
			n, err := io.Copy(file, r)
			if err != nil {
				return 0, err
			}
			if err := session.Wait(); err != nil {
				return 0, err
			}
			totalBytes += n
		}

	}
	return totalBytes, nil
}

func init() {
	distributors.Add("ssh", func() replica.Distributor {
		return &SSH{}
	})
}
