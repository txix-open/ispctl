package command

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"isp-ctl/command/common_config"
)

func CommonConfig() *cli.Command {
	return &cli.Command{
		Name:         "common",
		Usage:        "access to common configs commands",
		Action:       cc.action,
		BashComplete: cc.bachComplete,
		Subcommands: []*cli.Command{
			common_config.Get(),
			common_config.Set(),
			common_config.Delete(),
			common_config.Remove(),
			common_config.Link(),
			common_config.UnLink(),
			common_config.Contain(),
		},
	}
}

var cc commonConfigCommand

type commonConfigCommand struct{}

func (g commonConfigCommand) action(ctx *cli.Context) error {
	if ctx.Args().First() == "" {
		fmt.Print("need use command: ")
		for _, comm := range CommonConfig().Subcommands {
			fmt.Printf("[%s] ", comm.Name)
		}
		fmt.Printf("\n")
	}
	return nil
}

func (g commonConfigCommand) bachComplete(ctx *cli.Context) {
	for _, command := range CommonConfig().Subcommands {
		fmt.Println(command.Name)
	}
}
