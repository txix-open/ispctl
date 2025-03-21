package bash

import (
	"fmt"
	"ispctl/model"
	"ispctl/service"

	"github.com/urfave/cli/v2"

	"github.com/txix-open/bellows"
)

const (
	CommonConfigName = "config_name"
	ModuleName       = "module_name"
	ModuleData       = "module_data"

	VariableName = "variable_name"
	Empty        = "empty"
)

type BashService interface {
	GetAllVariables() ([]model.Variable, error)
	GetAvailableConfigs() ([]model.ModuleInfo, error)
	GetConfigurationByModuleName(moduleName string) (*model.Config, error)
}

type BeforeComplete func(ctx *cli.Context) (service.Config, error)

type Autocomplete struct {
	beforeComplete BeforeComplete
	service        BashService
}

func NewAutocomplete(beforeComplete BeforeComplete) Autocomplete {
	return Autocomplete{
		beforeComplete: beforeComplete,
	}
}

func (c Autocomplete) Complete(first string, second string) cli.BashCompleteFunc {
	return func(ctx *cli.Context) {
		cfgService, err := c.beforeComplete(ctx)
		if err != nil {
			return
		}
		c.service = cfgService
		c.complete(ctx, first, second)
	}
}

func (c Autocomplete) complete(ctx *cli.Context, first string, second string) {
	switch ctx.NArg() {
	case 0:
		c.completeFirstArg(ctx, first)
	case 1:
		c.completeSecondArg(ctx, second)
	}
}

func (c Autocomplete) completeFirstArg(ctx *cli.Context, first string) {
	switch first {
	case ModuleName:
		c.printAvailableModules()
	case VariableName:
		c.printAvailableVariables()
	}
}

func (c Autocomplete) completeSecondArg(ctx *cli.Context, second string) {
	switch second {
	case ModuleName:
		c.printAvailableModules()
	case ModuleData:
		c.printModuleData(ctx.Args().First())
	}
}

func (c Autocomplete) printAvailableVariables() {
	vars, err := c.service.GetAllVariables()
	if err != nil {
		return
	}
	for _, v := range vars {
		fmt.Println(v.Name)
	}
}

func (c Autocomplete) printAvailableModules() {
	modules, err := c.service.GetAvailableConfigs()
	if err != nil {
		return
	}
	for _, module := range modules {
		fmt.Println(module.Name)
	}
}

func (c Autocomplete) printModuleData(moduleName string) {
	config, err := c.service.GetConfigurationByModuleName(moduleName)
	if err != nil {
		return
	}
	for key := range bellows.Flatten(config.Data) {
		fmt.Printf(".%v\n", key)
	}
}
