package common_config

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"isp-ctl/bash"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
)

func Get() *cli.Command {
	return &cli.Command{
		Name:         "get",
		Usage:        "get common configurations",
		Action:       get.action,
		BashComplete: bash.Get(bash.CommonConfigName, bash.CommonConfigData).Complete,
	}
}

var get getCommand

type getCommand struct{}

func (g getCommand) action(ctx *cli.Context) error {
	if err := flag.CheckGlobal(ctx); err != nil {
		return err
	}

	configName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)

	commonConfigs, err := service.Config.GetMapNameCommonConfig()
	if err != nil {
		return err
	}

	if configName == "" {
		fmt.Printf("available next common configs:\n")
		for name := range commonConfigs {
			fmt.Printf("[%s] ", name)
		}
		fmt.Println()
		return nil
	}

	pathObject, err = utils.CheckPath(pathObject)
	if err != nil {
		return err
	}

	config, ok := commonConfigs[configName]
	if !ok {
		return errors.Errorf("common config [%s] not found", configName)
	}

	if pathObject == "" {
		utils.PrintAnswer(config.Data)
	} else if jsonObject, err := json.Marshal(config.Data); err != nil {
		return err
	} else {
		utils.CheckObject(jsonObject, pathObject)
	}
	return nil
}
