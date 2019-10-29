package commonConfig

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/pkg/errors"
	"isp-ctl/bash"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
)

func Get() cli.Command {
	return cli.Command{
		Name:         "get",
		Usage:        "get common configurations",
		Action:       get.action,
		BashComplete: bash.CommonConfig.GetSetDelete,
	}
}

var get getCommand

type getCommand struct{}

func (g getCommand) action(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		utils.PrintError(err)
		return
	}

	ccName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)

	commonConfigs, err := service.Config.GetMapCommonConfigByName()
	if err != nil {
		utils.PrintError(err)
		return
	}

	if ccName == "" {
		for name := range commonConfigs {
			fmt.Printf("%s ", name)
		}
		fmt.Println()
		return
	}

	pathObject, err = utils.CheckPath(pathObject)
	if err != nil {
		utils.PrintError(err)
		return
	}

	config, ok := commonConfigs[ccName]
	if !ok {
		utils.PrintError(errors.Errorf("common config %s not found", ccName))
		return
	}

	if pathObject == "" {
		utils.PrintAnswer(config.Data)
	} else if jsonObject, err := json.Marshal(config.Data); err != nil {
		utils.PrintError(err)
	} else {
		utils.CheckObject(jsonObject, pathObject)
	}
}
