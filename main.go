package main

import (
	"fmt"
	"log"
	"os"

	"ispctl/command"
	"ispctl/command/flag"

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
	app.Commands = command.AllCommands()
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
