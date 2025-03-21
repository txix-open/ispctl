package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"ispctl/bash"
	"ispctl/command"
	"ispctl/command/flag"
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
		flag.Host,
		flag.Unsafe,
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
	enableUnsafe := ctx.Bool(flag.Unsafe.Name)

	host := ctx.String(flag.Host.Name)
	host = strings.Replace(host, "'", "", -1)

	configCli, err := repository.NewGrpcClientWithHost(host)
	if err != nil {
		return service.Config{}, err
	}
	configRepo := repository.NewConfig(configCli)
	return service.NewConfig(enableUnsafe, configRepo), nil
}
