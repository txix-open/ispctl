package command

import (
	"ispctl/model"

	"github.com/txix-open/isp-kit/rc/schema"
	"github.com/urfave/cli/v2"
)

type AutoComplete interface {
	Complete(first string, second string) cli.BashCompleteFunc
}

type ConfigService interface {
	CreateUpdateConfig(stringToChange string, configuration *model.Config) (map[string]any, error)
	CreateUpdateConfigV2(configuration *model.Config) (map[string]any, error)
	DeleteVariable(variableName string) error
	GetAllVariables() ([]model.Variable, error)
	GetAvailableConfigs() ([]model.ModuleInfo, error)
	GetConfigurationByModuleName(moduleName string) (*model.Config, error)
	GetSchemaByModuleId(moduleId string) (schema.Schema, error)
	GetVariableByName(variableName string) (*model.Variable, error)
	UpsertVariables(vars []model.UpsertVariableRequest) error
}

func AllCommands(configService ConfigService, autoComplete AutoComplete) []*cli.Command {
	status := NewStatus(configService)
	get := NewGet(configService, autoComplete)
	set := NewSet(configService, autoComplete)
	delete := NewDelete(configService, autoComplete)
	schema := NewSchema(configService, autoComplete)
	merge := NewMerge(configService, autoComplete)
	gitGet := NewGitGet()
	variables := NewVariables(configService, autoComplete)
	return []*cli.Command{
		status.Command(),
		get.Command(),
		set.Command(),
		delete.Command(),
		schema.Command(),
		merge.Command(),
		gitGet.Command(),
		variables.Command(),
	}
}
