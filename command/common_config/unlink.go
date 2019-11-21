package common_config

import (
	"github.com/codegangsta/cli"
	"isp-ctl/bash"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
)

func UnLink() cli.Command {
	return cli.Command{
		Name:         "unlink",
		Usage:        "unlink common configurations from module configuration",
		Action:       unlink.action,
		BashComplete: bash.Get(bash.CommonConfigName, bash.ModuleName).Complete,
	}
}

var unlink unlinkCommand

type unlinkCommand struct{}

func (g unlinkCommand) action(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		utils.PrintError(err)
		return
	}

	configName := ctx.Args().First()
	moduleName := ctx.Args().Get(1)

	if err := service.Config.UnlinkCommonConfigFromModule(configName, moduleName); err != nil {
		utils.PrintError(err)
	}
}
