package command

import (
	"encoding/json"

	"github.com/urfave/cli/v2"
	"ispctl/bash"
	"ispctl/command/flag"
	"ispctl/command/utils"
	"ispctl/service"
)

func Get() *cli.Command {
	return &cli.Command{
		Name:         "get",
		Usage:        "get configuration by module_name",
		Action:       get.action,
		BashComplete: bash.Get(bash.ModuleName, bash.ModuleData).Complete,
	}
}

var get getCommand

type getCommand struct{}

func (g getCommand) action(ctx *cli.Context) error {
	if err := flag.CheckGlobal(ctx); err != nil {
		return err
	}
	moduleName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)

	config, err := service.Config.GetConfigurationByModuleName(moduleName)
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
