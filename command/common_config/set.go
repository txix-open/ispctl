package common_config

import (
	"encoding/json"
	"github.com/codegangsta/cli"
	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
	"isp-ctl/bash"
	"isp-ctl/cfg"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
)

func Set() cli.Command {
	return cli.Command{
		Name:         "set",
		Usage:        "set common configurations",
		Action:       set.action,
		BashComplete: bash.Get(bash.CommonConfigName, bash.CommonConfigData).Complete,
	}
}

var set setCommand

type setCommand struct{}

func (g setCommand) action(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		utils.PrintError(err)
		return
	}

	configName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)
	changeObject := ctx.Args().Get(2)

	if configName == "" {
		utils.PrintError(errors.New("empty config name"))
		return
	}

	mapConfigByName, err := service.Config.GetMapNameCommonConfig()
	if err != nil {
		utils.PrintError(err)
		return
	}

	config, ok := mapConfigByName[configName]
	if !ok {
		config = cfg.CommonConfig{
			Name: configName,
		}
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
		utils.CreateUpdateCommonConfig(changeObject, config)
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
			utils.CreateUpdateCommonConfig(stringToChange, config)
		}
	}
}
