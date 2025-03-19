package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"ispctl/command"
	"ispctl/flag"
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
	app.Commands = []*cli.Command{
		command.Status(),
		command.Get(),
		command.Set(),
		command.Delete(),
		command.Schema(),
		command.Merge(),
		command.GitGet(),
		command.VariablesCommands(),
	}
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
