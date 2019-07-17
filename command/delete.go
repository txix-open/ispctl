package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/tidwall/sjson"
)

func Delete() cli.Command {
	return cli.Command{
		Name:   "delete",
		Usage:  "delete configuration info by module_name",
		Action: deleteComm.action,
	}
}

var deleteComm deleteCommand

type deleteCommand struct{}

func (d deleteCommand) action(c *cli.Context) {
	if err := checkFlags(c); err != nil {
		fmt.Println(err)
		return
	}
	moduleName := c.Args().First()
	pathObject := c.Args().Get(1)

	moduleConfiguration, jsonObject := getModuleConfiguration(moduleName)
	if moduleConfiguration == nil {
		return
	}

	if pathObject == "" {
		createUpdateConfig("", moduleConfiguration)
	} else {
		if stringToChange, err := sjson.Delete(string(jsonObject), pathObject); err != nil {
			fmt.Println(err)
			return
		} else {
			createUpdateConfig(stringToChange, moduleConfiguration)
		}
	}
}
