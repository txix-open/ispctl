package bash

import (
	"fmt"
	"github.com/integration-system/bellows"
	"github.com/urfave/cli/v2"
	"isp-ctl/flag"
	"isp-ctl/service"
)

const (
	CommonConfigName bashArg = "config_name"
	CommonConfigData bashArg = "config_data"
	ModuleName       bashArg = "module_name"
	ModuleData       bashArg = "module_data"
	Empty            bashArg = "empty"
)

type bashArg string

type bashCommand struct {
	first  bashArg
	second bashArg
}

func Get(first bashArg, second bashArg) *bashCommand {
	return &bashCommand{
		first:  first,
		second: second,
	}
}

func (c bashCommand) Complete(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		return
	}
	commonConfigs, err := service.Config.GetMapNameCommonConfig()
	if err != nil {
		return
	}
	switch ctx.NArg() {
	case 0:
		switch c.first {
		case ModuleName:
			if arrayOfModules, err := service.Config.GetAvailableConfigs(); err != nil {
				return
			} else {
				for _, module := range arrayOfModules {
					fmt.Println(module.Name)
				}
			}
		case CommonConfigName:
			for _, config := range commonConfigs {
				fmt.Println(config.Name)
			}
		}
	case 1:
		switch c.second {
		case CommonConfigData:
			if config, ok := commonConfigs[ctx.Args().First()]; ok {
				for key, _ := range bellows.Flatten(config.Data) {
					fmt.Printf(".%v\n", key)
				}
			}
		case ModuleName:
			if arrayOfModules, err := service.Config.GetAvailableConfigs(); err != nil {
				return
			} else {
				for _, module := range arrayOfModules {
					fmt.Println(module.Name)
				}
			}
		case ModuleData:
			if moduleConfiguration, err := service.Config.GetConfigurationByModuleName(ctx.Args().First()); err != nil {
				return
			} else {
				for key, _ := range bellows.Flatten(moduleConfiguration.Data) {
					fmt.Printf(".%v\n", key)
				}
			}
		}
	}
}
