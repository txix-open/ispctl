package cfg

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/grpc/client"
)

type ConfigClient struct {
	cli     *client.Client
	address string
}

func NewConfigClient() *ConfigClient {
	return &ConfigClient{}
}

func (c *ConfigClient) ReceiveConfig(address string) error {
	cli, err := client.Default()
	if err != nil {
		return err
	}
	c.cli = cli
	cli.Upgrade([]string{address})

	c.address = address
	return nil
}

func (c *ConfigClient) GetAvailableConfigs() ([]ModuleInfo, error) {
	response := make([]ModuleInfo, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.cli.Invoke(getAvailableConfigs).
		JsonResponseBody(&response).
		Do(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getAvailableConfigs)
	}
	return response, nil
}

func (c *ConfigClient) GetConfigByModuleName(name string) (*Config, error) {
	request := &getModuleByUuidAndNameRequest{ModuleName: name}
	response := new(Config)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.cli.Invoke(getConfigByModuleName).
		JsonRequestBody(request).
		JsonResponseBody(&response).
		Do(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getConfigByModuleName)
	}
	return response, nil
}

func (c *ConfigClient) CreateUpdateConfig(request Config) (*Config, error) {
	response := new(Config)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.cli.Invoke(createUpdateConfig).
		JsonRequestBody(request).
		JsonResponseBody(&response).
		Do(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", createUpdateConfig)
	}
	return response, nil
}

func (c *ConfigClient) GetSchemaByModuleId(moduleId string) (*ConfigSchema, error) {
	request := getSchemaByModuleIdRequest{ModuleId: moduleId}
	response := new(ConfigSchema)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.cli.Invoke(getSchemaByModuleId).
		JsonRequestBody(request).
		JsonResponseBody(&response).
		Do(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getSchemaByModuleId)
	}
	return response, nil
}

func (c *ConfigClient) GetAllVariables() ([]Variable, error) {
	response := make([]Variable, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.cli.Invoke(getAllVariables).
		JsonResponseBody(&response).
		Do(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getAllVariables)
	}
	return response, nil
}

func (c *ConfigClient) GetVariableByName(variableName string) (*Variable, error) {
	request := VariableByNameRequest{Name: variableName}
	response := new(Variable)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.cli.Invoke(getVariableByName).
		JsonRequestBody(request).
		JsonResponseBody(&response).
		Do(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getVariableByName)
	}
	return response, nil
}

func (c *ConfigClient) UpsertVariables(request []UpsertVariableRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.cli.Invoke(upsertVariable).
		JsonRequestBody(request).
		Do(ctx)
	if err != nil {
		return errors.WithMessagef(err, "call %s", upsertVariable)
	}
	return nil
}

func (c *ConfigClient) DeleteVariable(variableName string) error {
	request := VariableByNameRequest{Name: variableName}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.cli.Invoke(deleteVariable).
		JsonRequestBody(request).
		Do(ctx)
	if err != nil {
		return errors.WithMessagef(err, "call %s", deleteVariable)
	}
	return nil
}
