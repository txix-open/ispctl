package command

import (
	"encoding/json"

	"ispctl/bash"
	"ispctl/command/utils"
	"ispctl/model"

	"github.com/tidwall/sjson"
	"github.com/urfave/cli/v2"
)

type SetService interface {
	GetConfigurationByModuleName(moduleName string) (*model.Config, error)
	CreateUpdateConfig(stringToChange string, configuration *model.Config) (map[string]any, error)
}

type Set struct {
	service      SetService
	autoComplete AutoComplete
}

func NewSet(service SetService, autoComplete AutoComplete) Set {
	return Set{
		service:      service,
		autoComplete: autoComplete,
	}
}

func (c Set) Command() *cli.Command {
	return &cli.Command{
		Name:         "set",
		Usage:        "set configuration by module_name",
		Action:       c.action,
		BashComplete: c.autoComplete.Complete(bash.ModuleName, bash.ModuleData),
	}
}

func (c Set) action(ctx *cli.Context) error {
	moduleName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)
	changeObject := ctx.Args().Get(2)

	config, err := c.service.GetConfigurationByModuleName(moduleName)
	if err != nil {
		return err
	}

	pathObject, err = utils.CheckPath(pathObject)
	if err != nil {
		return err
	}

	changeObject, err = utils.CheckChangeObject(changeObject)
	if err != nil {
		return err
	}

	if pathObject == "" {
		return CreateUpdateConfig(changeObject, config, c.service)
	}
	jsonObject, err := json.Marshal(config.Data)
	if err != nil {
		return err
	}

	changeArgument := utils.ParseSetObject(changeObject)
	stringToChange, err := sjson.Set(string(jsonObject), pathObject, changeArgument)
	if err != nil {
		return err
	}
	return CreateUpdateConfig(stringToChange, config, c.service)
}
