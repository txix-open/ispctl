package command

import (
	"github.com/codegangsta/cli"
	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"isp-ctl/bash"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
	"os"
)

func Set() cli.Command {
	return cli.Command{
		Name:         "set",
		Usage:        "set configuration by module_name",
		Action:       set.action,
		BashComplete: bash.ModuleNameAndConfigurationPath.Complete,
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

	if changeObject == "" {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			utils.PrintError(err)
			return
		}
		changeObject = string(bytes)
	}
	if changeObject == "" {
		utils.PrintError(errors.New("expected argument"))
		return
	}

	if pathObject == "" {
		utils.CreateUpdateConfig(changeObject, moduleConfiguration)
		return
	} else {
		changeArgument := utils.ParseSetObject(changeObject)
		if stringToChange, err := sjson.Set(string(jsonObject), pathObject, changeArgument); err != nil {
			utils.PrintError(err)
			return
		} else {
			utils.CreateUpdateConfig(stringToChange, moduleConfiguration)
			return
		}
	}
}
