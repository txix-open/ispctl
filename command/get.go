package command

import (
	"encoding/json"

	"ispctl/bash"
	"ispctl/command/utils"
	"ispctl/model"

	"github.com/urfave/cli/v2"
)

type GetService interface {
	GetConfigurationByModuleName(moduleName string) (*model.Config, error)
}

type Get struct {
	service      SetService
	autoComplete AutoComplete
}

func NewGet(service SetService, autoComplete AutoComplete) Get {
	return Get{
		service:      service,
		autoComplete: autoComplete,
	}
}

func (c Get) Command() *cli.Command {
	return &cli.Command{
		Name:         "get",
		Usage:        "get configuration by module_name",
		Action:       c.action,
		BashComplete: c.autoComplete.Complete(bash.ModuleName, bash.ModuleData),
	}
}

func (c Get) action(ctx *cli.Context) error {
	moduleName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)

	config, err := c.service.GetConfigurationByModuleName(moduleName)
	if err != nil {
		return err
	}

	data := config.Data

	pathObject, err = utils.CheckPath(pathObject)
	if err != nil {
		return err
	}

	if pathObject == "" {
		utils.PrintAnswer(data)
	} else {
		jsonObject, err := json.Marshal(data)
		if err != nil {
			return err
		}
		utils.CheckObject(jsonObject, pathObject)
	}
	return nil
}
