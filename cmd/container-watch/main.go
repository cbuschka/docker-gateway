package main

import (
	"fmt"
	"github.com/cbuschka/docker-gateway/internal/operator"
	"os"
)

func main() {
	err := operator.Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed: %s\n", err.Error())
		os.Exit(1)
	}
}
