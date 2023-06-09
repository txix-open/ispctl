package common_config

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
	"github.com/urfave/cli/v2"
	"ispctl/bash"
	"ispctl/command/utils"
	"ispctl/flag"
	"ispctl/service"
)

func Delete() *cli.Command {
	return &cli.Command{
		Name:         "delete",
		Usage:        "set common configurations",
		Action:       deleteComm.action,
		BashComplete: bash.Get(bash.CommonConfigName, bash.CommonConfigData).Complete,
	}
}

var deleteComm deleteCommand

type deleteCommand struct{}

func (g deleteCommand) action(ctx *cli.Context) error {
	if err := flag.CheckGlobal(ctx); err != nil {
		return err
	}

	configName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)

	if configName == "" {
		return errors.New("empty config name")
	}

	pathObject, err := utils.CheckPath(pathObject)
	if err != nil {
		return err
	}

	config, err := service.Config.GetCommonConfigByName(configName)
	if err != nil {
		return err
	}

	if pathObject == "" {
		return utils.CreateUpdateCommonConfig("", config)
	} else {
		jsonObject, err := json.Marshal(config.Data)
		if err != nil {
			return err
		}
		if stringToChange, err := sjson.Delete(string(jsonObject), pathObject); err != nil {
			return err
		} else {
			return utils.CreateUpdateCommonConfig(stringToChange, config)
		}
	}
}
