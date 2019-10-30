package command

import (
	"encoding/json"
	"github.com/codegangsta/cli"
	"github.com/tidwall/sjson"
	"isp-ctl/bash"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
)

func Set() cli.Command {
	return cli.Command{
		Name:         "set",
		Usage:        "set configuration by module_name",
		Action:       set.action,
		BashComplete: bash.Get(bash.ModuleName, bash.ModuleData).Complete,
	}
}

var set setCommand

type setCommand struct{}

func (s setCommand) action(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		utils.PrintError(err)
		return
	}

	moduleName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)
	changeObject := ctx.Args().Get(2)

	config, err := service.Config.GetConfigurationByModuleName(moduleName)
	if err != nil {
		utils.PrintError(err)
		return
	}

	pathObject, err = utils.CheckPath(pathObject)
	if err != nil {
		utils.PrintError(err)
		return
	}

	changeObject, err = utils.CheckChangeObject(changeObject)
	if err != nil {
		utils.PrintError(err)
		return
	}

	if pathObject == "" {
		utils.CreateUpdateConfig(changeObject, config)
	} else {
		jsonObject, err := json.Marshal(config.Data)
		if err != nil {
			utils.PrintError(err)
			return
		}

		changeArgument := utils.ParseSetObject(changeObject)
		if stringToChange, err := sjson.Set(string(jsonObject), pathObject, changeArgument); err != nil {
			utils.PrintError(err)
		} else {
			utils.CreateUpdateConfig(stringToChange, config)
		}
	}
}
