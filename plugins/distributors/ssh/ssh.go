package ssh

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"strings"

	"github.com/pkg/sftp"
	"github.com/rogercoll/replica"
	"github.com/rogercoll/replica/plugins/distributors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type SSH struct {
	Hosts        []string `toml:"hosts"`
	Passwords    []string `toml:"passwords"`
	Destinations []string `toml:"destinations"`
}

const sampleConfig = `
  ## Replicate files via SSH
  # This is expected to be the list of ssh config made of hostIP, port and user to do the replication, eg: "192.168.1.90:22:bob"
  hosts = []
  # Password to use with the corresponding config, if empty will use the loaded ssh keys with the system SSH Agent
  usernames = []
  # Backup remote destinations of the corresponding config
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

func parsePath(path string) string {
	return strings.Replace(path, "//", "/", -1)
}

func (s *SSH) Save(files []replica.Backup) (int64, error) {
	var totalBytes int64
	if len(s.Hosts) != len(s.Destinations) {
		return 0, errors.New("Invalid configuration")
	}
	for i, c := range s.Hosts {
		config := strings.Split(c, ":")
		if len(config) != 3 {
			return 0, errors.New("Invalid host format")
		}
		auths := []ssh.AuthMethod{}
		if s.Passwords[i] == "" {
			if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
				auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers))
				defer sshAgent.Close()
			}
		} else {
			auths = append(auths, ssh.Password(s.Passwords[i]))
		}
		clientConfig := &ssh.ClientConfig{
			User:            config[2],
			Auth:            auths,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		client, err := ssh.Dial("tcp", config[0]+":"+config[1], clientConfig)
		if err != nil {
			return 0, err
		}
		// open an SFTP session over an existing ssh connection.
		sftp, err := sftp.NewClient(client)
		if err != nil {
			return 0, err
		}
		defer sftp.Close()
		for _, fileName := range files {
			// Create the destination file
			dstFile, err := sftp.Create(parsePath(s.Destinations[i] + path.Base(fileName.Name())))
			if err != nil {
				return 0, err
			}
			defer dstFile.Close()

			// write to file
			n, err := dstFile.ReadFrom(fileName.Data())
			if err != nil {
				fmt.Println("heldfjasjfkldas")
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
