package utils

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/txix-open/isp-kit/json"
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
			pathObject = value
			continue
		}
		pathObject = fmt.Sprintf("%s.%s", pathObject, value)
	}
	return pathObject, nil
}

func CheckChangeObject(changeObject string) (string, error) {
	if changeObject == "" {
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		changeObject = string(bytes)
	}
	if changeObject == "" {
		return "", errors.New("expected argument")
	}
	return changeObject, nil
}

func PrintAnswer(data any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "	")
	return encoder.Encode(data)
}

func ParseSetObject(argument string) any {
	tryParse := []byte(argument)
	if tryParse[0] == '"' && tryParse[len(tryParse)-1] == '"' {
		tryParse = tryParse[1 : len(tryParse)-1]
		return string(tryParse)
	}

	if argument == "null" {
		return nil
	}

	mapStringInterface := make(map[string]any)
	if err := json.Unmarshal(tryParse, &mapStringInterface); err == nil {
		return mapStringInterface
	}

	arrayOfObject := make([]any, 0)
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

func CheckObject(jsonObject []byte, depth string) error {
	jsonString := gjson.GetBytes(jsonObject, depth)
	if jsonString.Raw == "" {
		return errors.Errorf("path '%s' not found\n", depth)
	}
	var data any
	err := json.Unmarshal([]byte(jsonString.Raw), &data)
	if err != nil {
		return err
	}
	return PrintAnswer(data)
}
