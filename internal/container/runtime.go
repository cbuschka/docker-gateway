package container

import (
	"context"
	"time"
)

type ContainerRuntime interface {
	ListContainers(ctx context.Context) ([]Container, error)
	Subscribe(ctx context.Context) (<-chan Event, <-chan error)
	LastModified() (time.Time, error)
}
