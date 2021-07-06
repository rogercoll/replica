# Replica

Replica is a tool to generate directory backups and replicate them in multiple environments. It is a plugin-driven application, and for the moment it distinguishes 2 plugin types:

1. [Backup Plugins](#backup-plugins) generate backups of a system
2. [Distributor Plugins](#distributor-plugins) distributes the previous backups among diferent systems

New plugins are designed to be easy to contribute, pull requests are super welcomed. For example, databases backup plugins or cloud services distributor like AWS S3.

## Installation:
### From Source:

Replica requires Go version 1.14 or newer, the Makefile requires GNU make.

1. [Install Go](https://golang.org/doc/install) >=1.14 (1.15 recommended)
2. Clone the Replica repository:
   ```
   git clone https://github.com/rogercoll/replica.git
   ```
3. Run `make` from the source directory
   ```
   make build
   sudo cp ./replica /usr/bin
   ```

## Backup Plugins

* [targz](./plugins/backup/targz)

## Distributor Plugins

* [local](./plugins/distributors/local)
* [ssh](./plugins/distributors/ssh)

## Configuration

Replica will try to find a default config file at these locations (in order):
1. $REPLICA_CONFIG_PATH
2. $HOME/.replica/replica.conf
3. /etc/replica/replica.conf

Example file:

```toml
[[backup.targz]]
  paths = ["/home/neck/hello.txt"]
  # Time format to be used in the backup file names (https://golang.org/pkg/time/#Time)
  # Default is layoutISO = "2006-01-02"

[[distributor.local]]
  destinations = ["/home/neck/backups/", "/home/neck/backups/test1/", "/home/neck/backups/test2/"]

[[distributor.ssh]]
  hosts = ["192.168.1.90:22:pi"]
  passwords = [""]
  destinations = ["/home/pi/extBackups/"]
```

**The program workflow is based on [Telegraf](https://github.com/influxdata/telegraf) agent, kudos to them for the great work!**
