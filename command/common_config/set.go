package common_config

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
	"github.com/urfave/cli/v2"
	"ispctl/bash"
	"ispctl/cfg"
	"ispctl/command/utils"
	"ispctl/flag"
	"ispctl/service"
)

func Set() *cli.Command {
	return &cli.Command{
		Name:         "set",
		Usage:        "set common configurations",
		Action:       set.action,
		BashComplete: bash.Get(bash.CommonConfigName, bash.CommonConfigData).Complete,
	}
}

var set setCommand

type setCommand struct{}

func (g setCommand) action(ctx *cli.Context) error {
	if err := flag.CheckGlobal(ctx); err != nil {
		return err
	}

	configName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)
	changeObject := ctx.Args().Get(2)

	if configName == "" {
		return errors.New("empty config name")
	}

	mapConfigByName, err := service.Config.GetMapNameCommonConfig()
	if err != nil {
		return err
	}

	config, ok := mapConfigByName[configName]
	if !ok {
		config = cfg.CommonConfig{
			Name: configName,
		}
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
		return utils.CreateUpdateCommonConfig(changeObject, config)
	} else {
		jsonObject, err := json.Marshal(config.Data)
		if err != nil {
			return err
		}
		changeArgument := utils.ParseSetObject(changeObject)
		if stringToChange, err := sjson.Set(string(jsonObject), pathObject, changeArgument); err != nil {
			return err
		} else {
			return utils.CreateUpdateCommonConfig(stringToChange, config)
		}
	}
}
