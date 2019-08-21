package command

import (
	"github.com/codegangsta/cli"
	"github.com/tidwall/sjson"
)

func Delete() cli.Command {
	return cli.Command{
		Name:         "delete",
		Usage:        "delete configuration by module_name",
		Action:       deleteComm.action,
		BashComplete: bash.run,
	}
}

var deleteComm deleteCommand

type deleteCommand struct{}

func (d deleteCommand) action(c *cli.Context) {
	if err := checkFlags(c); err != nil {
		printError(err)
		return
	}
	moduleName := c.Args().First()
	pathObject := c.Args().Get(1)

	moduleConfiguration, jsonObject := getModuleConfiguration(moduleName)
	if moduleConfiguration == nil {
		return
	}

	pathObject, ok := checkPath(pathObject)
	if !ok {
		return
	}

	if pathObject == "" {
		createUpdateConfig("", moduleConfiguration)
	} else {
		if stringToChange, err := sjson.Delete(string(jsonObject), pathObject); err != nil {
			printError(err)
			return
		} else {
			createUpdateConfig(stringToChange, moduleConfiguration)
		}
	}
}
