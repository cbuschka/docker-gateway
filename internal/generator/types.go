package generator

type DomainData struct {
	Key                    string
	Name                   string
	Hosts                  []HostData
	SslEnabled             bool
	SslCertificateFile     string
	SslKeyFile             string
	SslCertifcateChainFile string
}

type HostData struct {
	ContainerName string
	Ip            string
	Port          int
}

type ConfigData struct {
	Domains []DomainData
}
