package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"isp-ctl/command/commonConfig"
)

func CommonConfig() cli.Command {
	return cli.Command{
		Name:         "common",
		Usage:        "access to common configs commands",
		Action:       cc.action,
		BashComplete: cc.bachComplete,
		Subcommands: []cli.Command{
			commonConfig.Get(),
			commonConfig.Set(),
			commonConfig.Delete(),
			commonConfig.Remove(),
			commonConfig.Link(),
			commonConfig.UnLink(),
			commonConfig.Contain(),
		},
	}
}

var cc commonConfigCommand

type commonConfigCommand struct{}

func (g commonConfigCommand) action(ctx *cli.Context) {
	if ctx.Args().First() == "" {
		fmt.Print("need use command: ")
		for _, comm := range CommonConfig().Subcommands {
			fmt.Printf("[%s] ", comm.Name)
		}
		fmt.Printf("\n")
	}
}

func (g commonConfigCommand) bachComplete(ctx *cli.Context) {
	for _, command := range CommonConfig().Subcommands {
		fmt.Println(command.Name)
	}
}
