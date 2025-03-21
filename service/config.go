package service

import (
	"encoding/json"

	"ispctl/model"
	"ispctl/repository"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/rc/schema"
)

type Config struct {
	enableUnsafe bool
	configRepo   repository.Config
}

func NewConfig(enableUnsafe bool, configRepo repository.Config) Config {
	return Config{
		enableUnsafe: enableUnsafe,
		configRepo:   configRepo,
	}
}

func (c Config) GetAvailableConfigs() ([]model.ModuleInfo, error) {
	return c.configRepo.GetAvailableConfigs()
}

func (c Config) GetConfigurationByModuleName(moduleName string) (*model.Config, error) {
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

func (c Config) GetSchemaByModuleId(moduleId string) (schema.Schema, error) {
	configSchema, err := c.configRepo.GetSchemaByModuleId(moduleId)
	if err != nil {
		return nil, err
	}
	return configSchema.Schema, nil
}

func (c Config) CreateUpdateConfig(stringToChange string, configuration *model.Config) (map[string]any, error) {
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

func (c Config) CreateUpdateConfigV2(configuration *model.Config) (map[string]any, error) {
	configuration.Unsafe = c.enableUnsafe
	resp, err := c.configRepo.CreateUpdateConfig(*configuration)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c Config) GetAllVariables() ([]model.Variable, error) {
	return c.configRepo.GetAllVariables()
}

func (c Config) GetVariableByName(variableName string) (*model.Variable, error) {
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

func (c Config) DeleteVariable(variableName string) error {
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

func (c Config) UpsertVariables(vars []model.UpsertVariableRequest) error {
	return c.configRepo.UpsertVariables(vars)
}
