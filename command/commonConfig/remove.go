package commonConfig

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/pkg/errors"
	"isp-ctl/bash"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
)

func Remove() cli.Command {
	return cli.Command{
		Name:         "remove",
		Usage:        "remove common configurations",
		Action:       remove.action,
		BashComplete: bash.CommonConfig.Remove,
	}
}

var remove removeCommand

type removeCommand struct{}

func (g removeCommand) action(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		utils.PrintError(err)
		return
	}

	ccName := ctx.Args().First()

	if ccName == "" {
		utils.PrintError(errors.New("empty common config name"))
		return
	}

	commonConfigs, err := service.Config.GetMapCommonConfigByName()
	if err != nil {
		utils.PrintError(err)
		return
	}

	config, ok := commonConfigs[ccName]
	if !ok {
		utils.PrintError(errors.Errorf("common config %s not found", ccName))
		return
	}

	if deleted, err := service.Config.DeleteCommonConfig(config.Id); err != nil {
		utils.PrintError(err)
	} else {
		fmt.Printf("count deleted config: %d\n", deleted)
	}
}
