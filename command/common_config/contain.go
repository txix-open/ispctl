package common_config

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/pkg/errors"
	"isp-ctl/bash"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
)

func Contain() cli.Command {
	return cli.Command{
		Name:         "contain",
		Usage:        "availability common configuration in module",
		Action:       contain.action,
		BashComplete: bash.Get(bash.CommonConfigName, bash.Empty).Complete,
	}
}

var contain containCommand

type containCommand struct{}

func (g containCommand) action(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		utils.PrintError(err)
		return
	}

	configName := ctx.Args().First()

	if configName == "" {
		utils.PrintError(errors.New("empty common config name"))
		return
	}

	config, err := service.Config.GetCommonConfigByName(configName)
	if err != nil {
		utils.PrintError(err)
		return
	}

	links, err := service.Config.GetLinksCommonConfig(config.Id)
	if err != nil {
		utils.PrintError(err)
	} else {
		fmt.Printf("config [%s] linked with next modules:\n", configName)
		for _, link := range links {
			fmt.Printf("[%s] ", link)
		}
		fmt.Printf("\n")
	}
}
