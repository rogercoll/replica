package controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/rogercoll/replica/config"
)

// Controller runs a set of plugins.
type Controller struct {
	Config *config.Config
}

// NewController returns an Controller for the given Config.
func NewController(config *config.Config) (*Controller, error) {
	a := &Controller{
		Config: config,
	}
	return a, nil
}

func (c *Controller) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	//all backupfilters must be different to prevent data race
	fmt.Println(len(c.Config.Backups))
	backupFiles := make(map[string][]*os.File, len(c.Config.Backups))
	for _, running := range c.Config.Backups {
		//create logger with with filed of the backup name
		wg.Add(1)
		//contoller could have already consumed them
		go func() {
			defer wg.Done()
			files, err := running.Backup.Do()
			if err != nil {
				log.Println(err)
			}
			backupFiles[running.Name] = files
		}()
	}
	wg.Wait()
	for _, running := range c.Config.Auths {
		fmt.Println(running)
		//create logger with with filed of the backup name
		wg.Add(1)
		//contoller could have already consumed them
		go func() {
			defer wg.Done()
			for k, files := range backupFiles {
				fmt.Println("helkjfjejfe")
				sBytes, err := running.Auth.Save(files)
				if err != nil {
					log.Println(err)
				}
				fmt.Printf("Backup: [%s] Saved bytes[%s]\n", k, sBytes)
			}
		}()
	}
	wg.Wait()
	return nil
}
