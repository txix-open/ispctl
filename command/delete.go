package command

import (
	"encoding/json"

	"ispctl/bash"
	"ispctl/command/flag"
	"ispctl/command/utils"
	"ispctl/service"

	"github.com/tidwall/sjson"
	"github.com/urfave/cli/v2"
)

func Delete() *cli.Command {
	return &cli.Command{
		Name:         "delete",
		Usage:        "delete configuration by module_name",
		Before:       flag.ApplyGlobalFlags,
		Action:       deleteComm.action,
		BashComplete: bash.Get(bash.ModuleName, bash.ModuleData).Complete,
	}
}

var deleteComm deleteCommand

type deleteCommand struct{}

func (d deleteCommand) action(ctx *cli.Context) error {
	moduleName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)

	moduleConfiguration, err := service.Config.GetConfigurationByModuleName(moduleName)
	if err != nil {
		return err
	}

	pathObject, err = utils.CheckPath(pathObject)
	if err != nil {
		return err
	}

	if pathObject == "" {
		return utils.CreateUpdateConfig("", moduleConfiguration)
	} else {
		jsonObject, err := json.Marshal(moduleConfiguration.Data)
		if err != nil {
			return err
		}

		if stringToChange, err := sjson.Delete(string(jsonObject), pathObject); err != nil {
			return err
		} else {
			return utils.CreateUpdateConfig(stringToChange, moduleConfiguration)
		}
	}
}
