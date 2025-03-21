package command

import (
	"encoding/json"
	"os"

	"ispctl/bash"
	"ispctl/command/utils"
	"ispctl/model"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type MergeService interface {
	GetConfigurationByModuleName(moduleName string) (*model.Config, error)
	CreateUpdateConfigV2(configuration *model.Config) (map[string]any, error)
}

type Merge struct {
	service      MergeService
	autoComplete AutoComplete
}

func NewMerge(service MergeService, autoComplete AutoComplete) Merge {
	return Merge{
		service:      service,
		autoComplete: autoComplete,
	}
}

func (c Merge) Command() *cli.Command {
	return &cli.Command{
		Name:         "merge",
		Usage:        "merge actual config with config from stdin",
		BashComplete: c.autoComplete.Complete(bash.CommonConfigName, bash.Empty),
		Action:       c.action,
	}
}

func (c Merge) action(ctx *cli.Context) error {
	moduleName := ctx.Args().First()
	config, err := c.service.GetConfigurationByModuleName(moduleName)
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

	result, err := c.service.CreateUpdateConfigV2(config)
	if err != nil {
		return errors.WithMessagef(err, "update module config '%s'", moduleName)
	}
	utils.PrintAnswer(result)

	return nil
}
