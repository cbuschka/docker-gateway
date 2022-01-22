package operator

import (
	"context"
	"github.com/cbuschka/docker-gateway/internal/generator"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func (w *Watcher) handle() error {

	containers, err := w.runtime.ListContainers(context.Background())
	if err != nil {
		return err
	}

	domainsByKey := make(map[string]generator.DomainData, 0)
	for _, container := range containers {
		domainKey := toDomainKey(container.DomainName)
		domain, domainFound := domainsByKey[domainKey]
		if !domainFound {
			domain = generator.DomainData{Key: domainKey,
				Name:                   container.DomainName,
				Hosts:                  []generator.HostData{},
				SslEnabled:             true,
				SslCertificateFile:     "/etc/ssl/certs/selfsigned.crt",
				SslKeyFile:             "/etc/ssl/private/selfsigned.key",
				SslCertifcateChainFile: "/etc/ssl/certs/selfsigned.crt"}
		}

		host := generator.HostData{ContainerName: container.Name,
			Ip:   container.Ip,
			Port: container.Port}
		hosts := append(domain.Hosts, host)
		domain.Hosts = hosts
		domainsByKey[domainKey] = domain
	}

	domains := make([]generator.DomainData, 0)
	for _, domain := range domainsByKey {
		domains = append(domains, domain)
	}

	//FIXME check host reachability

	config, err := generator.Generate(w.config.TemplateFile, domains)
	if err != nil {
		return err
	}

	log.Printf("Nginx config generated:\n%s", config)

	err = w.writeNginxConfig(config)
	if err != nil {
		return err
	}

	err = w.restartNginx()
	if err != nil {
		return err
	}

	return nil
}

func (w *Watcher) restartNginx() error {

	commandLine := strings.Split(w.config.ReloadCommand, " ")

	command := exec.Command(commandLine[0], commandLine[1:]...)
	output, err := command.CombinedOutput()
	if err != nil {
		return err
	}

	log.Printf("Nginx restarted: exit code=%d, output=%s", command.ProcessState.ExitCode(), string(output))
	return nil
}

func toDomainKey(domainName string) string {

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(domainName, "_")
}

func (w *Watcher) writeNginxConfig(config []byte) error {

	configpath := w.config.OutputFile

	err := ioutil.WriteFile(configpath, config, 0644)
	if err != nil {
		return err
	}

	log.Printf("Nginx config written to %s.", configpath)

	return nil
}
