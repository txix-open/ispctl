package common_config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"isp-ctl/bash"
	"isp-ctl/flag"
	"isp-ctl/service"
)

func Contain() *cli.Command {
	return &cli.Command{
		Name:         "contain",
		Usage:        "availability common configuration in module",
		Action:       contain.action,
		BashComplete: bash.Get(bash.CommonConfigName, bash.Empty).Complete,
	}
}

var contain containCommand

type containCommand struct{}

func (g containCommand) action(ctx *cli.Context) error {
	if err := flag.CheckGlobal(ctx); err != nil {
		return err
	}

	configName := ctx.Args().First()

	if configName == "" {
		return errors.New("empty common config name")
	}

	config, err := service.Config.GetCommonConfigByName(configName)
	if err != nil {
		return err
	}

	links, err := service.Config.GetLinksCommonConfig(config.Id)
	if err != nil {
		return err
	} else {
		fmt.Printf("config [%s] linked with next modules:\n", configName)
		for _, link := range links {
			fmt.Printf("[%s] ", link)
		}
		fmt.Printf("\n")
	}
	return nil
}
