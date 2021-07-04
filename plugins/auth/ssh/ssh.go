package ssh

import (
	"os"

	"github.com/rogercoll/replica"
	"github.com/rogercoll/replica/plugins/auth"
)

type SSH struct {
	Hosts     []string `toml:"hosts"`
	Usernames []string `toml:"usernames"`
}

const sampleConfig = `
  ## Replicate files via SSH
  # This is expected to be the list of hosts to do the replication
  hosts = []
  # Username for the ssh connection for the previously defined hosts
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

func (s *SSH) Save([]*os.File) (int64, error) {
	return 0, nil
}

func init() {
	auth.Add("ssh", func() replica.Authenticator {
		return &SSH{}
	})
}
