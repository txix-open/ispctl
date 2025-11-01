package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"ispctl/bash"
	"ispctl/command"
	"ispctl/repository"
	"ispctl/service"

	"github.com/urfave/cli/v2"
)

var version = "1.1.0"

func main() {
	initCommands()
}

func initCommands() {
	app := cli.NewApp()
	app.Usage = "isp configurations updater"

	app.Version = version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     command.HostFlagName,
			Usage:    command.HostFlagUsage,
			Value:    "127.0.0.1:9002",
			Required: false,
			Aliases:  []string{"g", "configAddr"},
		},
		&cli.BoolFlag{
			Name:  command.UnsafeFlagName,
			Usage: command.UnsafeFlagUsage,
		},
	}
	app.EnableBashCompletion = true

	var configService service.Config
	var autoComplete bash.Autocomplete
	app.Before = func(ctx *cli.Context) error {
		service, err := configServiceFromGlobalFlags(ctx)
		if err != nil {
			return err
		}
		configService = service
		return nil
	}
	autoComplete = bash.NewAutocomplete(configServiceFromGlobalFlags)
	app.Commands = command.AllCommands(&configService, autoComplete)

	app.BashComplete = func(context *cli.Context) {
		for _, appCommand := range app.Commands {
			fmt.Println(appCommand.Name)
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func configServiceFromGlobalFlags(ctx *cli.Context) (service.Config, error) {
	enableUnsafe := ctx.Bool(command.UnsafeFlagName)

	host := ctx.String(command.HostFlagName)
	host = strings.ReplaceAll(host, "'", "")

	configCli, err := repository.NewGrpcClientWithHost(host)
	if err != nil {
		return service.Config{}, err
	}
	configRepo := repository.NewConfig(configCli)
	return service.NewConfig(enableUnsafe, configRepo), nil
}
