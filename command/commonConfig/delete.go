package commonConfig

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/integration-system/bellows"
	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
)

func Delete() cli.Command {
	return cli.Command{
		Name:         "delete",
		Usage:        "set common configurations",
		Action:       deleteComm.action,
		BashComplete: deleteComm.bashComplete,
	}
}

var deleteComm deleteCommand

type deleteCommand struct{}

func (g deleteCommand) action(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		utils.PrintError(err)
		return
	}

	ccName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)

	if ccName == "" {
		utils.PrintError(errors.New("empty config name"))
		return
	}

	pathObject, err := utils.CheckPath(pathObject)
	if err != nil {
		utils.PrintError(err)
		return
	}

	commonConfigs, err := service.Config.GetMapCommonConfigByName()
	if err != nil {
		utils.PrintError(err)
		return
	}

	config, ok := commonConfigs[ccName]
	if !ok {
		utils.PrintError(errors.Errorf("common config %s not found", ccName))
		return
	}

	if pathObject == "" {
		utils.CreateUpdateCommonConfig("", config)
	} else {
		jsonObject, err := json.Marshal(config.Data)
		if err != nil {
			utils.PrintError(err)
			return
		}
		if stringToChange, err := sjson.Delete(string(jsonObject), pathObject); err != nil {
			utils.PrintError(err)
		} else {
			utils.CreateUpdateCommonConfig(stringToChange, config)
		}
	}
}

func (g deleteCommand) bashComplete(ctx *cli.Context) {
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
