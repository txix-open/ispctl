package cfg

import (
	"context"
	"time"

	"github.com/integration-system/isp-kit/grpc/client"
	"github.com/pkg/errors"
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
		ReadJsonResponse(&response).
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
		ReadJsonResponse(&response).
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
		ReadJsonResponse(&response).
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
		ReadJsonResponse(&response).
		Do(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getSchemaByModuleId)
	}
	return response, nil
}

func (c *ConfigClient) GetCommonConfigs(req []string) ([]CommonConfig, error) {
	response := make([]CommonConfig, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.cli.Invoke(getCommonConfigs).
		JsonRequestBody(req).
		ReadJsonResponse(&response).
		Do(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getCommonConfigs)
	}
	return response, nil
}

func (c *ConfigClient) CreateUpdateCommonConfig(req CommonConfig) (*CommonConfig, error) {
	response := new(CommonConfig)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.cli.Invoke(createUpdateCommonConfig).
		JsonRequestBody(req).
		ReadJsonResponse(&response).
		Do(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", createUpdateCommonConfig)
	}
	return response, nil
}

func (c *ConfigClient) DeleteCommonConfig(id string) (*Deleted, error) {
	request := identityRequest{Id: id}
	response := new(Deleted)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.cli.Invoke(deleteCommonConfig).
		JsonRequestBody(request).
		ReadJsonResponse(&response).
		Do(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", deleteCommonConfig)
	}
	return response, nil
}

func (c *ConfigClient) GetLinksCommonConfig(id string) (*Links, error) {
	request := identityRequest{Id: id}
	response := new(Links)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.cli.Invoke(getLinksCommonConfig).
		JsonRequestBody(request).
		ReadJsonResponse(&response).
		Do(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getLinksCommonConfig)
	}
	return response, nil
}

func (c *ConfigClient) CompileCommonConfigs(request CompileConfigs) (map[string]interface{}, error) {
	response := make(map[string]interface{})
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.cli.Invoke(compileCommonConfigs).
		JsonRequestBody(request).
		ReadJsonResponse(response).
		Do(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", compileCommonConfigs)
	}
	return response, nil
}

func (c *ConfigClient) errorHandler(err error) error {
	return errors.WithMessagef(err, "isp-config-service: %s", c.address)
}
