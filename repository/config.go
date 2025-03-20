package repository

import (
	"context"
	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/grpc/client"
	"ispctl/entity"
	"time"
)

const (
	getAvailableConfigs   = "config/module/get_modules_info"
	getConfigByModuleName = "config/config/get_active_config_by_module_name"
	createUpdateConfig    = "config/config/create_update_config"
	getSchemaByModuleId   = "config/schema/get_by_module_id"

	getAllVariables   = "config/variable/all"
	getVariableByName = "config/variable/get_by_name"
	upsertVariable    = "config/variable/upsert"
	deleteVariable    = "config/variable/delete"
)

type Config struct {
	cli     *client.Client
	baseCtx context.Context
}

func NewConfig(cli *client.Client) Config {
	return Config{
		cli:     cli,
		baseCtx: context.Background(),
	}
}

func (c Config) GetAvailableConfigs() ([]entity.ModuleInfo, error) {
	response := make([]entity.ModuleInfo, 0)
	err := c.cli.Invoke(getAvailableConfigs).
		JsonResponseBody(&response).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getAvailableConfigs)
	}
	return response, nil
}

func (c Config) GetConfigByModuleName(name string) (*entity.Config, error) {
	request := &entity.GetModuleByUuidAndNameRequest{ModuleName: name}
	response := new(entity.Config)
	err := c.cli.Invoke(getConfigByModuleName).
		JsonRequestBody(request).
		JsonResponseBody(&response).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getConfigByModuleName)
	}
	return response, nil
}

func (c Config) CreateUpdateConfig(request entity.Config) (*entity.Config, error) {
	response := new(entity.Config)
	err := c.cli.Invoke(createUpdateConfig).
		JsonRequestBody(request).
		JsonResponseBody(&response).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", createUpdateConfig)
	}
	return response, nil
}

func (c Config) GetSchemaByModuleId(moduleId string) (*entity.ConfigSchema, error) {
	request := entity.GetSchemaByModuleIdRequest{ModuleId: moduleId}
	response := new(entity.ConfigSchema)
	err := c.cli.Invoke(getSchemaByModuleId).
		JsonRequestBody(request).
		JsonResponseBody(&response).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getSchemaByModuleId)
	}
	return response, nil
}

func (c Config) GetAllVariables() ([]entity.Variable, error) {
	response := make([]entity.Variable, 0)
	err := c.cli.Invoke(getAllVariables).
		JsonResponseBody(&response).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getAllVariables)
	}
	return response, nil
}

func (c Config) GetVariableByName(variableName string) (*entity.Variable, error) {
	request := entity.VariableByNameRequest{Name: variableName}
	response := new(entity.Variable)
	err := c.cli.Invoke(getVariableByName).
		JsonRequestBody(request).
		JsonResponseBody(&response).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getVariableByName)
	}
	return response, nil
}

func (c Config) UpsertVariables(request []entity.UpsertVariableRequest) error {
	err := c.cli.Invoke(upsertVariable).
		JsonRequestBody(request).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return errors.WithMessagef(err, "call %s", upsertVariable)
	}
	return nil
}

func (c Config) DeleteVariable(variableName string) error {
	request := entity.VariableByNameRequest{Name: variableName}
	err := c.cli.Invoke(deleteVariable).
		JsonRequestBody(request).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return errors.WithMessagef(err, "call %s", deleteVariable)
	}
	return nil
}
