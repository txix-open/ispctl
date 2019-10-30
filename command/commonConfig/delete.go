package commonConfig

import (
	"encoding/json"
	"github.com/codegangsta/cli"
	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
	"isp-ctl/bash"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
)

func Delete() cli.Command {
	return cli.Command{
		Name:         "delete",
		Usage:        "set common configurations",
		Action:       deleteComm.action,
		BashComplete: bash.CommonConfig.ConfigName_ConfigData,
	}
}

var deleteComm deleteCommand

type deleteCommand struct{}

func (g deleteCommand) action(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		utils.PrintError(err)
		return
	}

	ccName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)

	if ccName == "" {
		utils.PrintError(errors.New("empty config name"))
		return
	}

	pathObject, err := utils.CheckPath(pathObject)
	if err != nil {
		utils.PrintError(err)
		return
	}

	commonConfigs, err := service.Config.GetMapCommonConfigByName()
	if err != nil {
		utils.PrintError(err)
		return
	}

	config, ok := commonConfigs[ccName]
	if !ok {
		utils.PrintError(errors.Errorf("common config [%s] not found", ccName))
		return
	}

	if pathObject == "" {
		utils.CreateUpdateCommonConfig("", config)
	} else {
		jsonObject, err := json.Marshal(config.Data)
		if err != nil {
			utils.PrintError(err)
			return
		}
		if stringToChange, err := sjson.Delete(string(jsonObject), pathObject); err != nil {
			utils.PrintError(err)
		} else {
			utils.CreateUpdateCommonConfig(stringToChange, config)
		}
	}
}
