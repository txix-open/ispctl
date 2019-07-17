package command

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/tidwall/sjson"
	"os"
	"strings"
)

func Set() cli.Command {
	return cli.Command{
		Name:   "set",
		Usage:  "set configuration info",
		Action: set.action,
	}
}

var set setCommand

type setCommand struct{}

func (s setCommand) action(c *cli.Context) {
	if err := checkFlags(c); err != nil {
		fmt.Println(err)
		return
	}

	moduleName := c.Args().First()
	pathObject := c.Args().Get(1)
	changeObject := c.Args().Get(2)

	moduleConfiguration, jsonObject := getModuleConfiguration(moduleName)
	if moduleConfiguration == nil {
		return
	}

	if changeObject == "" {
		fmt.Print("Enter new value: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		changeObject = scanner.Text()
	}

	if changeObject == "" {
		fmt.Println("Error. Expected argument")
		return
	}

	if pathObject == "" {
		createUpdateConfig(changeObject, moduleConfiguration)
		return
	} else {
		changeArgument := s.parseObject(changeObject)
		if stringToChange, err := sjson.Set(string(jsonObject), pathObject, changeArgument); err != nil {
			fmt.Println(err)
			return
		} else {
			createUpdateConfig(stringToChange, moduleConfiguration)
			return
		}
	}
}

func (s setCommand) parseObject(argument string) interface{} {
	argument = strings.Replace(argument, " ", "", -1)
	tryParse := []byte(argument)
	if tryParse[0] == '"' && tryParse[len(tryParse)-1] == '"' {
		tryParse = tryParse[1 : len(tryParse)-1]
		return string(tryParse)
	}

	if argument == "null" {
		return nil
	}

	mapStringInterface := make(map[string]interface{})
	if err := json.Unmarshal(tryParse, &mapStringInterface); err == nil {
		return mapStringInterface
	}

	arrayOfObject := make([]interface{}, 0)
	if err := json.Unmarshal(tryParse, &arrayOfObject); err == nil {
		return arrayOfObject
	}

	var integer int64
	if err := json.Unmarshal(tryParse, &integer); err == nil {
		return integer
	}

	var float float64
	if err := json.Unmarshal(tryParse, &float); err == nil {
		return float
	}

	var flag bool
	if err := json.Unmarshal(tryParse, &flag); err == nil {
		return flag
	}

	return argument
}
