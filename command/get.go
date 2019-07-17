package command

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/tidwall/gjson"
)

func Get() cli.Command {
	return cli.Command{
		Name:   "get",
		Usage:  "get configuration by module_name",
		Action: get.action,
	}
}

var get getCommand

type getCommand struct{}

func (g getCommand) action(c *cli.Context) {
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

	pathObject, ok := checkPath(pathObject)
	if !ok {
		return
	}

	if pathObject == "" {
		printAnswer(moduleConfiguration.Data)
	} else {
		g.checkObject(jsonObject, pathObject)
	}
}

func (g getCommand) checkObject(jsonObject []byte, depth string) {
	jsonString := gjson.Get(string(jsonObject), depth)
	if jsonString.Raw == "" {
		fmt.Printf("Path '%s' not found\n", depth)
	} else {
		var data interface{}
		if err := json.Unmarshal([]byte(jsonString.Raw), &data); err != nil {
			fmt.Println(err)
		} else {
			printAnswer(data)
		}
	}
}
