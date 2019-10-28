package command

import (
	"encoding/json"
	"github.com/codegangsta/cli"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"isp-ctl/bash"
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
		printAnswer(moduleConfiguration.Data)
	} else {
		g.checkObject(jsonObject, pathObject)
	}
}

func (g getCommand) checkObject(jsonObject []byte, depth string) {
	jsonString := gjson.Get(string(jsonObject), depth)
	if jsonString.Raw == "" {
		printError(errors.Errorf("Path '%s' not found\n", depth))
	} else {
		var data interface{}
		if err := json.Unmarshal([]byte(jsonString.Raw), &data); err != nil {
			printError(err)
		} else {
			printAnswer(data)
		}
	}
}
