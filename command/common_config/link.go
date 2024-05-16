package common_config

import (
	"github.com/urfave/cli/v2"
	"ispctl/bash"
	"ispctl/flag"
	"ispctl/service"
)

func Link() *cli.Command {
	return &cli.Command{
		Name:         "link",
		Usage:        "link common configurations to module configuration",
		Action:       link.action,
		BashComplete: bash.Get(bash.CommonConfigName, bash.ModuleName).Complete,
	}
}

var link linkCommand

type linkCommand struct{}

func (g linkCommand) action(ctx *cli.Context) error {
	if err := flag.CheckGlobal(ctx); err != nil {
		return err
	}

	configName := ctx.Args().First()
	moduleName := ctx.Args().Get(1)

	if err := service.Config.LinkCommonConfigToModule(configName, moduleName); err != nil {
		return err
	}
	return nil
}
