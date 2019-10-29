package commonConfig

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/pkg/errors"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
)

func Link() cli.Command {
	return cli.Command{
		Name:         "link",
		Usage:        "link common configurations to module configuration",
		Action:       link.action,
		BashComplete: link.bashComplete,
	}
}

var link linkCommand

type linkCommand struct{}

func (g linkCommand) action(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		utils.PrintError(err)
		return
	}

	configName := ctx.Args().First()
	moduleName := ctx.Args().Get(1)

	if configName == "" {
		utils.PrintError(errors.New("need config name"))
		return
	}

	if moduleName == "" {
		utils.PrintError(errors.New("need module name"))
		return
	}

	commonConfigs, err := service.Config.GetMapCommonConfigByName()
	if err != nil {
		utils.PrintError(err)
		return
	}

	config, ok := commonConfigs[configName]
	if !ok {
		utils.PrintError(errors.Errorf("common config %s not found", configName))
		return
	}

	moduleConfiguration, _, err := service.Config.GetConfigurationAndJsonByModuleName(moduleName)
	if err != nil {
		utils.PrintError(err)
		return
	}

	ccIdNameMap := make(map[string]string)
	for name, value := range commonConfigs {
		ccIdNameMap[value.Id] = name
	}

	needLink := true
	for _, configId := range moduleConfiguration.CommonConfigs {
		if config.Id == configId {
			needLink = false
			break
		}
	}

	if needLink {
		moduleConfiguration.CommonConfigs = append(moduleConfiguration.CommonConfigs, config.Id)
		if linked, err := service.Config.UpdateConfigAndGetLinkCommon(*moduleConfiguration); err != nil {
			utils.PrintError(err)
		} else {
			for _, name := range linked {
				fmt.Printf("[%s] ", ccIdNameMap[name])
			}
			fmt.Printf("\n")
		}
	} else {
		for _, name := range moduleConfiguration.CommonConfigs {
			fmt.Printf("[%s] ", ccIdNameMap[name])
		}
		fmt.Printf("\n")
	}
}

func (g linkCommand) bashComplete(ctx *cli.Context) {
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
