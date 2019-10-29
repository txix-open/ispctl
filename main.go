package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"isp-ctl/command"
	"isp-ctl/flag"
	"isp-ctl/service"
	"log"
	"os"
)

var version = "1.1.0"

func main() {
	defineConfigurationPath()
	initCommands()
}

func initCommands() {
	app := cli.NewApp()
	app.Usage = "isp configurations updater"
	app.UsageText =
		`	ispctl [flag...] status
	ispctl [flag...] get 		module_name  property_path
	ispctl [flag...] set 		module_name  property_path  [new_object]
	ispctl [flag...] delete 		module_name  property_path
	ispctl [flag...] schema		module_name  [local_flag]`

	app.Version = version
	app.Flags = []cli.Flag{
		flag.Host,
		flag.Uuid,
		flag.Color,
		flag.Unsafe,
	}
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		command.Status(),
		command.Get(),
		command.Set(),
		command.Delete(),
		command.Schema(),
		command.CommonConfig(),
	}
	app.BashComplete = func(context *cli.Context) {
		for _, command := range app.Commands {
			fmt.Println(command.Name)
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func defineConfigurationPath() {
	if os.Getenv("APP_MODE") == "dev" {
		service.SetConfigurationPath("./conf/config.yml")
	} else {
		service.SetConfigurationPath("/etc/ispctl/config.yml")
	}
}
