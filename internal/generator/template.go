package generator

import (
	"bytes"
	"io/ioutil"
	"text/template"
)

func Generate(templateFile string, domains []DomainData) ([]byte, error) {

	templateBytes, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}

	tmpl := template.New(templateFile)
	tmpl, err = tmpl.Parse(string(templateBytes))
	if err != nil {
		return nil, err
	}

	bytesBuf := bytes.NewBuffer([]byte{})

	err = tmpl.Execute(bytesBuf, ConfigData{Domains: domains})
	if err != nil {
		return nil, err
	}

	return bytesBuf.Bytes(), nil
}
