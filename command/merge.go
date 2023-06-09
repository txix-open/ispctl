package command

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"ispctl/bash"
	"ispctl/command/utils"
	"ispctl/flag"
	"ispctl/service"
)

func Merge() *cli.Command {
	return &cli.Command{
		Name:         "merge",
		Usage:        "merge actual config with config from stdin",
		BashComplete: bash.Get(bash.CommonConfigName, bash.Empty).Complete,
		Action: func(context *cli.Context) error {
			err := flag.CheckGlobal(context)
			if err != nil {
				return err
			}
			moduleName := context.Args().First()
			config, err := service.Config.GetConfigurationByModuleName(moduleName)
			if err != nil {
				return errors.WithMessagef(err, "get module config '%s'", moduleName)
			}

			inputConfig := make(map[string]any)
			err = json.NewDecoder(os.Stdin).Decode(&inputConfig)
			if err != nil {
				return errors.WithMessage(err, "json unmarshal input config")
			}
			for key, value := range inputConfig {
				config.Data[key] = value
			}

			result, err := service.Config.CreateUpdateConfigV2(config)
			if err != nil {
				return errors.WithMessagef(err, "update module config '%s'", moduleName)
			}
			utils.PrintAnswer(result)

			return nil
		},
	}
}
