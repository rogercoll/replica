package controller

import (
	"context"
	"log"
	"sync"

	"github.com/rogercoll/replica/config"
	"github.com/rogercoll/replica/models"
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
	backupFiles := make(map[string][]string, len(c.Config.Backups))
	for _, running := range c.Config.Backups {
		//create logger with with filed of the backup name
		wg.Add(1)
		//contoller could have already consumed them
		go func() {
			defer wg.Done()
			files, err := running.Backup.Do()
			if err != nil {
				log.Println(err)
			} else {
				backupFiles[running.Name] = files
			}
		}()
	}
	wg.Wait()
	for _, running := range c.Config.Distributors {
		//create logger with with filed of the backup name
		wg.Add(1)
		//contoller could have already consumed them
		go func(r *models.RunningDistributor) {
			defer wg.Done()
			for k, files := range backupFiles {
				sBytes, err := r.Dist.Save(files)
				if err != nil {
					r.Log.Error(err)
				} else {
					r.Log.Printf("Backup: [%s] Saved bytes: [%d]\n", k, sBytes)
				}
			}
		}(running)
	}
	wg.Wait()
	return nil
}
