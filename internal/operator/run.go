package operator

import configPkg "github.com/cbuschka/docker-gateway/internal/config"

func Run() error {
	config, err := configPkg.GetConfig()
	if err != nil {
		return err
	}

	w, err := NewWatcher(config)
	if err != nil {
		return err
	}

	return w.Run()
}
