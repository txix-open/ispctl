package common_config

import (
	"github.com/codegangsta/cli"
	"isp-ctl/bash"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
)

func Link() cli.Command {
	return cli.Command{
		Name:         "link",
		Usage:        "link common configurations to module configuration",
		Action:       link.action,
		BashComplete: bash.Get(bash.CommonConfigName, bash.ModuleName).Complete,
	}
}

var link linkCommand

type linkCommand struct{}

func (g linkCommand) action(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		utils.PrintError(err)
		return
	}

	configName := ctx.Args().First()
	moduleName := ctx.Args().Get(1)

	if err := service.Config.LinkCommonConfigToModule(configName, moduleName); err != nil {
		utils.PrintError(err)
	}
}
