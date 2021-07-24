package controller

import (
	"context"
	"log"
	"sync"

	"github.com/rogercoll/replica"
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
	backupFiles := make(map[string][]replica.Backup, len(c.Config.Backups))
	for _, running := range c.Config.Backups {
		//create logger with with filed of the backup name
		wg.Add(1)
		//contoller could have already consumed them
		go func() {
			defer wg.Done()
			results := make(chan replica.Backup)
			go running.Backup.Do(results)
			//read until channel is closed by the backup
			for result := range results {
				if result.Error() != nil {
					log.Println(result.Error())
				} else {
					backupFiles[running.Name] = append(backupFiles[running.Name], result)
				}
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
