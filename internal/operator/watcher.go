package operator

import (
	"context"
	"github.com/cbuschka/docker-gateway/internal/config"
	"github.com/cbuschka/docker-gateway/internal/container"
	"log"
	"time"
)

type Watcher struct {
	runtime container.ContainerRuntime
	config  config.Config
}

func NewWatcher(config config.Config) (*Watcher, error) {
	runtime, err := container.OpenDocker()
	if err != nil {
		return nil, err
	}
	return &Watcher{runtime: runtime, config: config}, nil
}

func (w *Watcher) Run() error {

	log.Printf("Generating initially...")
	err := w.handle()
	if err != nil {
		return err
	}

	modified := false
	lastModified := time.Now()
	evCh, errCh := w.runtime.Subscribe(context.Background())
	for {
		select {
		case ev := <-evCh:
			if ev.Type == "container die" || ev.Type == "container start" {
				log.Printf("Received docker event: %s", ev.Type)
				modified = true
			}
			break
		case <-time.Tick(1 * time.Second):
			if modified && lastModified.Add(5*time.Second).Before(time.Now()) {
				log.Printf("Regenerating nginx config...")
				err := w.handle()
				if err != nil {
					log.Printf("Failure: %v", err)
				} else {
					modified = false
					lastModified = time.Now()
				}
			}
		case err := <-errCh:
			return err
		}
	}
}
