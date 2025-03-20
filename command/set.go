package command

import (
	"encoding/json"

	"github.com/tidwall/sjson"
	"github.com/urfave/cli/v2"
	"ispctl/bash"
	"ispctl/command/flag"
	"ispctl/command/utils"
	"ispctl/service"
)

func Set() *cli.Command {
	return &cli.Command{
		Name:         "set",
		Usage:        "set configuration by module_name",
		Action:       set.action,
		BashComplete: bash.Get(bash.ModuleName, bash.ModuleData).Complete,
	}
}

var set setCommand

type setCommand struct{}

func (s setCommand) action(ctx *cli.Context) error {
	if err := flag.CheckGlobal(ctx); err != nil {
		return err
	}

	moduleName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)
	changeObject := ctx.Args().Get(2)

	config, err := service.Config.GetConfigurationByModuleName(moduleName)
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
		return utils.CreateUpdateConfig(changeObject, config)
	} else {
		jsonObject, err := json.Marshal(config.Data)
		if err != nil {
			return err
		}

		changeArgument := utils.ParseSetObject(changeObject)
		if stringToChange, err := sjson.Set(string(jsonObject), pathObject, changeArgument); err != nil {
			return err
		} else {
			return utils.CreateUpdateConfig(stringToChange, config)
		}
	}
}
