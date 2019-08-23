package command

import (
	"github.com/codegangsta/cli"
	"github.com/tidwall/sjson"
	"isp-ctl/bash"
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
		printError(err)
		return
	}
	moduleName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)

	moduleConfiguration, jsonObject, err := service.Config.GetConfigurationAndJsonByModuleName(moduleName)
	if err != nil {
		printError(err)
		return
	}

	pathObject, err = checkPath(pathObject)
	if err != nil {
		printError(err)
		return
	}

	if pathObject == "" {
		createUpdateConfig("", moduleConfiguration)
	} else {
		if stringToChange, err := sjson.Delete(string(jsonObject), pathObject); err != nil {
			printError(err)
			return
		} else {
			createUpdateConfig(stringToChange, moduleConfiguration)
		}
	}
}
