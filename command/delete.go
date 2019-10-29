package command

import (
	"github.com/codegangsta/cli"
	"github.com/tidwall/sjson"
	"isp-ctl/bash"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
)

func Delete() cli.Command {
	return cli.Command{
		Name:         "delete",
		Usage:        "delete configuration by module_name",
		Action:       deleteComm.action,
		BashComplete: bash.ModuleNameAndConfigurationPath.Complete,
	}
}

var deleteComm deleteCommand

type deleteCommand struct{}

func (d deleteCommand) action(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		utils.PrintError(err)
		return
	}
	moduleName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)

	moduleConfiguration, jsonObject, err := service.Config.GetConfigurationAndJsonByModuleName(moduleName)
	if err != nil {
		utils.PrintError(err)
		return
	}

	pathObject, err = utils.CheckPath(pathObject)
	if err != nil {
		utils.PrintError(err)
		return
	}

	if pathObject == "" {
		utils.CreateUpdateConfig("", moduleConfiguration)
	} else {
		if stringToChange, err := sjson.Delete(string(jsonObject), pathObject); err != nil {
			utils.PrintError(err)
			return
		} else {
			utils.CreateUpdateConfig(stringToChange, moduleConfiguration)
		}
	}
}
