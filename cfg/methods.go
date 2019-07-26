package cfg

import (
	"github.com/integration-system/isp-lib/http"
)

const (
	instanceIdHeader = "x-instance-identity"
	scheme           = "http://"
)

func NewConfigClient(client http.RestClient) *configClient {
	return &configClient{
		client:  client,
		headers: make(map[string]string),
	}
}

type configClient struct {
	client  http.RestClient
	headers map[string]string

	gateHost     string
	instanceUuid string
}

func (c *configClient) ReceiveConfig(gateHost, instanceUuid string) {
	c.instanceUuid = instanceUuid
	c.headers[instanceIdHeader] = instanceUuid
	c.gateHost = scheme + gateHost
}

func (c *configClient) GetAvailableConfigs() ([]ModuleInfo, error) {
	response := make([]ModuleInfo, 0)
	if err := c.client.Invoke("POST", c.gateHost+getAvailableConfigs, c.headers, nil, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *configClient) GetConfigByModuleName(name string) (*Config, error) {
	request := &getModuleByUuidAndNameRequest{Name: name, Uuid: c.instanceUuid}
	response := new(Config)
	if err := c.client.Invoke("POST", c.gateHost+getConfigByModuleName, c.headers, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *configClient) CreateUpdateConfig(request Config) (*Config, error) {
	response := new(Config)
	if err := c.client.Invoke("POST", c.gateHost+createUpdateConfig, c.headers, &request, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *configClient) GetSchemaByModuleId(moduleId int32) (*ConfigSchema, error) {
	request := &getSchemaByModuleIdRequest{ModuleId: moduleId}
	response := new(ConfigSchema)
	if err := c.client.Invoke("POST", c.gateHost+getSchemaByModuleId, c.headers, &request, response); err != nil {
		return nil, err
	}
	return response, nil
}
