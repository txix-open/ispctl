package utils

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"isp-ctl/service"
	"os"
	"strings"
)

func CheckPath(pathObject string) (string, error) {
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

func CheckChangeObject(changeObject string) (string, error) {
	if changeObject == "" {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		changeObject = string(bytes)
	}
	if changeObject == "" {
		return "", errors.New("expected argument")
	} else {
		return changeObject, nil
	}
}

func PrintAnswer(data interface{}) {
	if answer, err := json.MarshalIndent(data, "", "    "); err != nil {
		PrintError(err)
	} else {
		service.ColorService.Print(answer)
	}
}

func PrintError(err error) {
	fmt.Println("ERROR:", err)
	os.Exit(-1)
}

func ParseSetObject(argument string) interface{} {
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

func CheckObject(jsonObject []byte, depth string) {
	jsonString := gjson.Get(string(jsonObject), depth)
	if jsonString.Raw == "" {
		PrintError(errors.Errorf("path '%s' not found\n", depth))
	} else {
		var data interface{}
		if err := json.Unmarshal([]byte(jsonString.Raw), &data); err != nil {
			PrintError(err)
		} else {
			PrintAnswer(data)
		}
	}
}
