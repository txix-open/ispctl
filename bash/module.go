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

func (module) ModuleName_ModuleData(ctx *cli.Context) {
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
		moduleConfiguration, err := service.Config.GetConfigurationByModuleName(ctx.Args().Get(1))
		if err != nil {
			return
		} else {
			for key, _ := range bellows.Flatten(moduleConfiguration.Data) {
				fmt.Printf(".%v\n", key)
			}
		}
	}
}

func (module) ModuleName(ctx *cli.Context) {
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
