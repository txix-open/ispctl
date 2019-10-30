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

func Contain() cli.Command {
	return cli.Command{
		Name:         "contain",
		Usage:        "availability common configuration in module",
		Action:       contain.action,
		BashComplete: bash.CommonConfig.ConfigName,
	}
}

var contain containCommand

type containCommand struct{}

func (g containCommand) action(ctx *cli.Context) {
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
		utils.PrintError(errors.Errorf("common config [%s] not found", ccName))
		return
	}

	links, err := service.Config.GetLinksCommonConfig(config.Id)
	if err != nil {
		utils.PrintError(err)
	} else {
		fmt.Printf("%v\n", links)
	}
}
