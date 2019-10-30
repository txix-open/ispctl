package command

import (
	"encoding/json"
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
		BashComplete: bash.Module.ModuleName_ModuleData,
		Flags: []cli.Flag{
			flag.WithCommonConfig,
		},
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

	config, err := service.Config.GetConfigurationByModuleName(moduleName)
	if err != nil {
		utils.PrintError(err)
		return
	}

	data := config.Data
	if ctx.Bool(flag.WithCommonConfig.Name) {
		if data, err = service.Config.CompileDataWithCommonConfigs(data, config.CommonConfigs); err != nil {
			utils.PrintError(err)
			return
		}
	}

	pathObject, err = utils.CheckPath(pathObject)
	if err != nil {
		utils.PrintError(err)
		return
	}

	if pathObject == "" {
		utils.PrintAnswer(data)
	} else {
		jsonObject, err := json.Marshal(data)
		if err != nil {
			utils.PrintError(err)
			return
		}
		utils.CheckObject(jsonObject, pathObject)
	}
}
