package bash

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/integration-system/bellows"
	"isp-ctl/flag"
	"isp-ctl/service"
)

var (
	ModuleNameAndConfigurationPath = bashAction{enableConfigurationBash: true}
	ModuleName                     = bashAction{enableConfigurationBash: false}
)

type (
	bashAction struct{ enableConfigurationBash bool }
)

func (b bashAction) Complete(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		return
	}
	switch ctx.NArg() {
	case 0:
		b.moduleNameHelper()
	case 1:
		if b.enableConfigurationBash {
			b.configurationPathHelper(ctx.Args().First(), ctx.Args().Get(1))
		}
	}
}

func (bashAction) moduleNameHelper() {
	if arrayOfModules, err := service.Config.GetAvailableConfigs(); err != nil {
		return
	} else {
		for _, module := range arrayOfModules {
			fmt.Println(module.Name)
		}
	}
}

func (bashAction) configurationPathHelper(moduleName, path string) {
	if moduleConfiguration, _, err := service.Config.GetConfigurationAndJsonByModuleName(moduleName); err != nil {
		return
	} else {
		for key, _ := range bellows.Flatten(moduleConfiguration.Data) {
			fmt.Printf(".%v\n", key)
		}
	}
}
