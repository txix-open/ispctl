package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/integration-system/bellows"
	"isp-ctl/service"
)

var (
	bash       = bashAction{enableConfigurationBash: true}
	bashSchema = bashAction{enableConfigurationBash: false}
)

type (
	bashAction struct{ enableConfigurationBash bool }
)

func (b bashAction) run(ctx *cli.Context) {
	if err := checkFlags(ctx); err != nil {
		return
	}
	switch ctx.NArg() {
	case 0:
		b.moduleNameHelper()
	case 1:
		if b.enableConfigurationBash {
			b.configurationHelper(ctx.Args().First(), ctx.Args().Get(1))
		}
	}
}

func (bashAction) moduleNameHelper() {
	if arrayOfModules, err := service.ConfigClient.GetAvailableConfigs(); err != nil {
		return
	} else {
		for _, module := range arrayOfModules {
			fmt.Println(module.Name)
		}
	}
}

func (bashAction) configurationHelper(moduleName, path string) {
	moduleConfiguration, _ := getModuleConfiguration(moduleName)
	if moduleConfiguration == nil {
		return
	}
	for key, _ := range bellows.Flatten(moduleConfiguration.Data) {
		fmt.Printf(".%v\n", key)
	}
}
