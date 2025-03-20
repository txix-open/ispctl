package bash

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"ispctl/command/flag"
	"ispctl/service"

	"github.com/txix-open/bellows"
)

const (
	CommonConfigName bashArg = "config_name"
	ModuleName       bashArg = "module_name"
	ModuleData       bashArg = "module_data"

	VariableName bashArg = "variable_name"
	Empty        bashArg = "empty"
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
	err := flag.ApplyGlobalFlags(ctx)
	if err != nil {
		return
	}

	switch ctx.NArg() {
	case 0:
		c.completeFirstArg(ctx)
	case 1:
		c.completeSecondArg(ctx)
	}
}

func (c bashCommand) completeFirstArg(ctx *cli.Context) {
	switch c.first {
	case ModuleName:
		printAvailableModules()
	case VariableName:
		printAvailableVariables()
	}
}

func (c bashCommand) completeSecondArg(ctx *cli.Context) {
	switch c.second {
	case ModuleName:
		printAvailableModules()
	case ModuleData:
		printModuleData(ctx.Args().First())
	}
}

func printAvailableVariables() {
	vars, err := service.Config.GetAllVariables()
	if err != nil {
		return
	}
	for _, v := range vars {
		fmt.Println(v.Name)
	}
}

func printAvailableModules() {
	modules, err := service.Config.GetAvailableConfigs()
	if err != nil {
		return
	}
	for _, module := range modules {
		fmt.Println(module.Name)
	}
}

func printModuleData(moduleName string) {
	config, err := service.Config.GetConfigurationByModuleName(moduleName)
	if err != nil {
		return
	}
	for key := range bellows.Flatten(config.Data) {
		fmt.Printf(".%v\n", key)
	}
}
