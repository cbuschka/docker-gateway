package container

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"strconv"
	"sync"
	"time"
)

type Docker struct {
	client       *client.Client
	lock         *sync.Mutex
	lastModified time.Time
}

func OpenDocker() (ContainerRuntime, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	return &Docker{client: dockerClient, lock: &sync.Mutex{}, lastModified: time.Unix(0, 0)}, nil
}

func (docker *Docker) Subscribe(ctx context.Context) (<-chan Event, <-chan error) {

	eventCh := make(chan Event)
	errCh := make(chan error)

	go func() {
		dockerEventCh, dockerErrCh := docker.client.Events(ctx, types.EventsOptions{})
		for {
			select {
			case ev := <-dockerEventCh:
				eventCh <- Event{fmt.Sprintf("%s %s", ev.Type, ev.Action)}
			case err := <-dockerErrCh:
				errCh <- err
			case <-ctx.Done():
				return
			}
		}
	}()

	return eventCh, errCh
}

func (docker *Docker) ListContainers(ctx context.Context) ([]Container, error) {
	containers := make([]Container, 0)
	dockerContainers, err := docker.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	for _, dockerContainer := range dockerContainers {
		labels := dockerContainer.Labels
		ingressDomain, ingressDomainFound := labels["ingress-domain"]
		if !ingressDomainFound {
			continue
		}

		port := 80
		portStr, portStrFound := labels["ingress-port"]
		if portStrFound {
			port, err = strconv.Atoi(portStr)
			if err != nil {
				return nil, err
			}
		}

		ip := ""
		for _, network := range dockerContainer.NetworkSettings.Networks {
			ip = network.IPAddress
		}

		container := Container{
			ID:         dockerContainer.ID,
			Name:       dockerContainer.Names[0],
			DomainName: ingressDomain,
			Port:       port,
			Ip:         ip,
		}
		containers = append(containers, container)
	}

	return containers, nil
}

func (docker *Docker) LastModified() (time.Time, error) {
	docker.lock.Lock()
	defer docker.lock.Unlock()
	return docker.lastModified, nil
}
