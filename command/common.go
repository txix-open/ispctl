package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"isp-ctl/cfg"
	"isp-ctl/service"
	"strings"
)

func checkPath(pathObject string) (string, error) {
	str := strings.Split(pathObject, ".")
	pathObject = ""
	if len(str) == 1 || str[0] != "" {
		return pathObject, errors.New("path must start with '.'")
	}
	for key, value := range str {
		if key == 0 {
			continue
		}
		if key == 1 {
			pathObject = fmt.Sprintf("%s", value)
			continue
		}
		pathObject = fmt.Sprintf("%s.%s", pathObject, value)
	}
	return pathObject, nil
}

func createUpdateConfig(stringToChange string, configuration *cfg.Config) {
	if answer, err := service.Config.CreateUpdateConfig(stringToChange, configuration); err != nil {
		printError(err)
	} else if answer != nil {
		printAnswer(answer)
	}
}

func printAnswer(data interface{}) {
	if answer, err := json.MarshalIndent(data, "", "    "); err != nil {
		printError(err)
	} else {
		service.ColorService.Print(answer)
	}
}

func printError(err error) {
	fmt.Println("ERROR:", err)
}
