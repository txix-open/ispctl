package bash

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/integration-system/bellows"
	"isp-ctl/flag"
	"isp-ctl/service"
)

var CommonConfig commonConfig

type commonConfig struct{}

func (commonConfig) ConfigName_ModuleName(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		return
	}
	commonConfigs, err := service.Config.GetMapCommonConfigByName()
	if err != nil {
		return
	}
	switch ctx.NArg() {
	case 0:
		for _, config := range commonConfigs {
			fmt.Println(config.Name)
		}
	case 1:
		if arrayOfModules, err := service.Config.GetAvailableConfigs(); err != nil {
			return
		} else {
			for _, module := range arrayOfModules {
				fmt.Println(module.Name)
			}
		}
	}
}

func (commonConfig) ConfigName_ConfigData(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		return
	}
	commonConfigs, err := service.Config.GetMapCommonConfigByName()
	if err != nil {
		return
	}
	switch ctx.NArg() {
	case 0:
		for _, config := range commonConfigs {
			fmt.Println(config.Name)
		}
	case 1:
		if config, ok := commonConfigs[ctx.Args().First()]; ok {
			for key, _ := range bellows.Flatten(config.Data) {
				fmt.Printf(".%v\n", key)
			}
		}
	}
}

func (commonConfig) ConfigName(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		return
	}
	commonConfigs, err := service.Config.GetMapCommonConfigByName()
	if err != nil {
		return
	}
	switch ctx.NArg() {
	case 0:
		for _, config := range commonConfigs {
			fmt.Println(config.Name)
		}
	}
}
