# Plugins

To check the documentation of any plugin you can run:

```bash
replica -bck-filter targz -sample-config
```

For example, to get the sample configuration of targz backup and ssh and local distributors we can run:

```bash
replica -bck-filter targz -dist-filter ssh:local -sample-config
```

And the corresponding output:

```toml
###############################################################################
#                           BACKUP PLUGINS                                    #
###############################################################################

# Generate Tar Gz files for given paths
[[backup.targz]]
  ## Create Tar Gz files for given paths (directories)
  # This is expected to be the list of absolute host path directories
  paths = []
  # Time format to be used in the backup file names (https://golang.org/pkg/time/#Time)
  # Default is layoutISO = "2006-01-02"
  timeformat = "2006-01-02T15:04:05Z07:00"


###############################################################################
#                            DISTRIBUTOR PLUGINS                              #
###############################################################################

# Connect via SSH and copy files via SCP
[[distributor.ssh]]
  ## Replicate files via SSH
  # This is expected to be the list of ssh config made of hostIP, port and user to do the replication, eg: "192.168.1.90:22:bob"
  hosts = []
  # Password to use with the corresponding config, if empty will use the loaded ssh keys
  usernames = []
  # Backup remote destinations of the corresponding config
  usernames = []


# Store backupfiles to a local directories
[[distributor.local]]
  ## Copy backup files to other directories
  # Absolute direcotries path to store the backup files
  # Default is $HOME/backups
  destination = []
```
