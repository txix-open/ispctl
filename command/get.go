package command

import (
	"github.com/codegangsta/cli"
	"isp-ctl/bash"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
)

func Get() cli.Command {
	return cli.Command{
		Name:         "get",
		Usage:        "get configuration by module_name",
		Action:       get.action,
		BashComplete: bash.ModuleNameAndConfigurationPath.Complete,
	}
}

var get getCommand

type getCommand struct{}

func (g getCommand) action(ctx *cli.Context) {
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
		utils.PrintAnswer(moduleConfiguration.Data)
	} else {
		utils.CheckObject(jsonObject, pathObject)
	}
}
