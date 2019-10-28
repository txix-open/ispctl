package command

import (
	"encoding/json"
	"github.com/codegangsta/cli"
	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"isp-ctl/bash"
	"isp-ctl/flag"
	"isp-ctl/service"
	"os"
)

func Set() cli.Command {
	return cli.Command{
		Name:         "set",
		Usage:        "set configuration by module_name",
		Action:       set.action,
		BashComplete: bash.ModuleNameAndConfigurationPath.Complete,
	}
}

var set setCommand

type setCommand struct{}

func (s setCommand) action(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		printError(err)
		return
	}

	moduleName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)
	changeObject := ctx.Args().Get(2)

	moduleConfiguration, jsonObject, err := service.Config.GetConfigurationAndJsonByModuleName(moduleName)
	if err != nil {
		printError(err)
		return
	}

	pathObject, err = checkPath(pathObject)
	if err != nil {
		printError(err)
		return
	}

	if changeObject == "" {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			printError(err)
			return
		}
		changeObject = string(bytes)
	}
	if changeObject == "" {
		printError(errors.New("expected argument"))
		return
	}

	if pathObject == "" {
		createUpdateConfig(changeObject, moduleConfiguration)
		return
	} else {
		changeArgument := s.parseObject(changeObject)
		if stringToChange, err := sjson.Set(string(jsonObject), pathObject, changeArgument); err != nil {
			printError(err)
			return
		} else {
			createUpdateConfig(stringToChange, moduleConfiguration)
			return
		}
	}
}

func (s setCommand) parseObject(argument string) interface{} {
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
