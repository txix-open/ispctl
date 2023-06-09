package common_config

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"ispctl/bash"
	"ispctl/flag"
	"ispctl/service"
)

func Remove() *cli.Command {
	return &cli.Command{
		Name:         "remove",
		Usage:        "remove common configurations",
		Action:       remove.action,
		BashComplete: bash.Get(bash.CommonConfigName, bash.Empty).Complete,
	}
}

var remove removeCommand

type removeCommand struct{}

func (g removeCommand) action(ctx *cli.Context) error {
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

	links, deleted, err := service.Config.DeleteCommonConfig(config.Id)
	if err != nil {
		return err
	}

	if deleted {
		fmt.Printf("config [%s] deleted\n", config.Name)
	} else {
		fmt.Printf("config [%s] not deleted, need unlink in next modules:\n%v\n", config.Name, links)
	}
	return nil
}
