package common_config

import (
	"github.com/urfave/cli/v2"
	"isp-ctl/bash"
	"isp-ctl/flag"
	"isp-ctl/service"
)

func UnLink() *cli.Command {
	return &cli.Command{
		Name:         "unlink",
		Usage:        "unlink common configurations from module configuration",
		Action:       unlink.action,
		BashComplete: bash.Get(bash.CommonConfigName, bash.ModuleName).Complete,
	}
}

var unlink unlinkCommand

type unlinkCommand struct{}

func (g unlinkCommand) action(ctx *cli.Context) error {
	if err := flag.CheckGlobal(ctx); err != nil {
		return err
	}

	configName := ctx.Args().First()
	moduleName := ctx.Args().Get(1)

	if err := service.Config.UnlinkCommonConfigFromModule(configName, moduleName); err != nil {
		return err
	}
	return nil
}
