package main

import (
	"github.com/codegangsta/cli"
	"isp-ctl/command"
	"isp-ctl/flag"
	"isp-ctl/service"
	"log"
	"os"
)

func main() {
	defineConfigurationPath()
	initCommands()
}

func initCommands() {
	app := cli.NewApp()
	app.Usage = "cfg updater"
	app.UsageText =
		`	ispctl status 		[flag...]
	ispctl get 		[flag...]	module_name [property_path]
	ispctl set 		[flag...]	module_name [property_path] [new_object]
	ispctl delete 		[flag...]	module_name [property_path]`

	app.Flags = []cli.Flag{
		flag.Host,
		flag.Uuid,
		flag.Color,
	}

	app.Commands = []cli.Command{
		command.Status(),
		command.Get(),
		command.Set(),
		command.Delete(),
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
