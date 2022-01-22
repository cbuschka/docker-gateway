package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	TemplateFile  string `yaml:"templateFile"`
	OutputFile    string `yaml:"outputFile"`
	ReloadCommand string `yaml:"reloadCommand"`
}

func GetConfig() (Config, error) {

	configFile := "/etc/container-watch/config.yml"
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return Config{}, err
	}

	log.Printf("Config loaded from %s:", configFile)
	log.Printf("\ttemplateFile: %s", config.TemplateFile)
	log.Printf("\toutputFile: %s", config.OutputFile)
	log.Printf("\treloadCommand: %s", config.ReloadCommand)

	return config, err
}
