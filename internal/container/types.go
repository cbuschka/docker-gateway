package container

type Container struct {
	ID         string
	Name       string
	DomainName string
	Ip         string
	Port       int
}

type Event struct {
	Type string
}
