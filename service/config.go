package service

import (
	"encoding/json"

	"ispctl/model"
	"ispctl/repository"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/rc/schema"
)

var Config = &ConfigService{}

type ConfigService struct {
	UnsafeEnable bool
	configRepo   repository.Config
}

func (c *ConfigService) ReceiveConfiguration(host string) error {
	cli, err := repository.NewClientFromHost(host)
	if err != nil {
		return err
	}
	c.configRepo = repository.NewConfig(cli)
	return nil
}

func (c *ConfigService) GetAvailableConfigs() ([]model.ModuleInfo, error) {
	return c.configRepo.GetAvailableConfigs()
}

func (c *ConfigService) GetConfigurationByModuleName(moduleName string) (*model.Config, error) {
	if moduleName == "" {
		return nil, errors.New("need module name")
	}

	moduleConfiguration, err := c.configRepo.GetConfigByModuleName(moduleName)
	switch {
	case errors.Is(err, model.ErrModuleNotFound):
		return nil, errors.Errorf("module [%s] not found", moduleName)
	default:
		return moduleConfiguration, err
	}
}

func (c *ConfigService) GetSchemaByModuleId(moduleId string) (schema.Schema, error) {
	configSchema, err := c.configRepo.GetSchemaByModuleId(moduleId)
	if err != nil {
		return nil, err
	}
	return configSchema.Schema, nil
}

func (c *ConfigService) CreateUpdateConfig(stringToChange string, configuration *model.Config) (map[string]any, error) {
	newData := make(map[string]any)
	if stringToChange != "" {
		err := json.Unmarshal([]byte(stringToChange), &newData)
		if err != nil {
			return nil, err
		}
	}

	configuration.Data = newData
	return c.CreateUpdateConfigV2(configuration)
}

func (c *ConfigService) CreateUpdateConfigV2(configuration *model.Config) (map[string]any, error) {
	configuration.Unsafe = c.UnsafeEnable
	resp, err := c.configRepo.CreateUpdateConfig(*configuration)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *ConfigService) GetAllVariables() ([]model.Variable, error) {
	return c.configRepo.GetAllVariables()
}

func (c *ConfigService) GetVariableByName(variableName string) (*model.Variable, error) {
	if variableName == "" {
		return nil, errors.New("need variable name")
	}

	variable, err := c.configRepo.GetVariableByName(variableName)
	switch {
	case errors.Is(err, model.ErrVariableNotFound):
		return nil, errors.Errorf("variable [%s] not found", variableName)
	case err != nil:
		return nil, err
	default:
		return variable, nil
	}
}

func (c *ConfigService) DeleteVariable(variableName string) error {
	if variableName == "" {
		return errors.New("need variable name")
	}

	err := c.configRepo.DeleteVariable(variableName)
	switch {
	case errors.Is(err, model.ErrVariableNotFound):
		return errors.Errorf("variable [%s] not found", variableName)
	default:
		return err
	}
}

func (c *ConfigService) UpsertVariables(vars []model.UpsertVariableRequest) error {
	return c.configRepo.UpsertVariables(vars)
}
