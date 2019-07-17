package service

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"isp-ctl/conf"
)

var YamlService yamlService

type yamlService struct{}

func (yamlService) Parse() (*conf.Configuration, error) {
	resp := new(conf.Configuration)
	reader, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	} else {
		err := yaml.Unmarshal(reader, resp)
		return resp, err
	}
}
