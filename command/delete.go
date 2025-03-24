package command

import (
	"encoding/json"

	"ispctl/bash"
	"ispctl/command/utils"
	"ispctl/model"

	"github.com/tidwall/sjson"
	"github.com/urfave/cli/v2"
)

type DeleteService interface {
	GetConfigurationByModuleName(moduleName string) (*model.Config, error)
	CreateUpdateConfig(stringToChange string, configuration *model.Config) (map[string]any, error)
}

type Delete struct {
	service      DeleteService
	autoComplete AutoComplete
}

func NewDelete(service DeleteService, autoComplete AutoComplete) Delete {
	return Delete{
		service:      service,
		autoComplete: autoComplete,
	}
}

func (c Delete) Command() *cli.Command {
	return &cli.Command{
		Name:         "delete",
		Usage:        "delete configuration by module_name",
		Action:       c.action,
		BashComplete: c.autoComplete.Complete(bash.ModuleName, bash.ModuleData),
	}
}

func (c Delete) action(ctx *cli.Context) error {
	moduleName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)

	moduleConfiguration, err := c.service.GetConfigurationByModuleName(moduleName)
	if err != nil {
		return err
	}

	pathObject, err = utils.CheckPath(pathObject)
	if err != nil {
		return err
	}
	if pathObject == "" {
		return CreateUpdateConfig("", moduleConfiguration, c.service)
	}

	jsonObject, err := json.Marshal(moduleConfiguration.Data)
	if err != nil {
		return err
	}
	stringToChange, err := sjson.Delete(string(jsonObject), pathObject)
	if err != nil {
		return err
	}

	return CreateUpdateConfig(stringToChange, moduleConfiguration, c.service)
}
