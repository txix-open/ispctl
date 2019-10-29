package bash

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/integration-system/bellows"
	"isp-ctl/flag"
	"isp-ctl/service"
)

var Module module

type module struct{}

func (module) GetSetDelete(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		return
	}
	switch ctx.NArg() {
	case 0:
		if arrayOfModules, err := service.Config.GetAvailableConfigs(); err != nil {
			return
		} else {
			for _, module := range arrayOfModules {
				fmt.Println(module.Name)
			}
		}
	case 1:
		if moduleConfiguration, _, err := service.Config.GetConfigurationAndJsonByModuleName(ctx.Args().Get(1)); err != nil {
			return
		} else {
			for key, _ := range bellows.Flatten(moduleConfiguration.Data) {
				fmt.Printf(".%v\n", key)
			}
		}
	}
}

func (module) Schema(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		return
	}
	switch ctx.NArg() {
	case 0:
		if arrayOfModules, err := service.Config.GetAvailableConfigs(); err != nil {
			return
		} else {
			for _, module := range arrayOfModules {
				fmt.Println(module.Name)
			}
		}
	}
}
