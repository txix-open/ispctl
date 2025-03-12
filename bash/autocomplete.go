package bash

import (
	"fmt"
	"github.com/urfave/cli/v2"

	"github.com/txix-open/bellows"
	"ispctl/flag"
	"ispctl/service"
)

const (
	CommonConfigName bashArg = "config_name"
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
		}
	case 1:
		switch c.second {
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
