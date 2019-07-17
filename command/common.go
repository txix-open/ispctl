package command

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/integration-system/isp-lib/http"
	"github.com/valyala/fasthttp"
	"isp-ctl/cfg"
	"isp-ctl/flag"
	"isp-ctl/service"
	"strings"
)

func checkFlags(c *cli.Context) error {
	var (
		uuid, host string
		color      bool
	)
	colorFlag := strings.Split(flag.ColorName, ", ")
	if c.GlobalBool(colorFlag[0]) != false {
		color = true
	} else {
		color = c.GlobalBool(colorFlag[1])
	}
	service.ColorService.Enable = color

	service.ColorService.Print(nil)

	hostFlag := strings.Split(flag.GateHostName, ", ")
	if c.GlobalString(hostFlag[0]) != "" {
		host = c.GlobalString(hostFlag[0])
	} else {
		host = c.GlobalString(hostFlag[1])
	}

	uuidFlag := strings.Split(flag.InstanceUuidName, ", ")
	if c.GlobalString(uuidFlag[0]) != "" {
		uuid = c.GlobalString(uuidFlag[0])
	} else {
		uuid = c.GlobalString(uuidFlag[1])
	}

	if uuid == "" || host == "" {
		if configuration, err := service.YamlService.Parse(); err != nil {
			return err
		} else {
			if uuid == "" {
				uuid = configuration.InstanceUuid
			}
			if host == "" {
				host = configuration.GateHost
			}
		}
	}
	service.ConfigService.ReceiveConfiguration(host, uuid)
	return nil
}

func checkPath(pathObject string) (string, bool) {
	str := strings.Split(pathObject, ".")
	pathObject = ""
	if len(str) == 1 || str[0] != "" {
		fmt.Println("Path must start with '.'")
		return pathObject, false
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
	return pathObject, true
}

func printAnswer(data interface{}) {
	if answer, err := json.MarshalIndent(data, "", "    "); err != nil {
		fmt.Println(err)
	} else {
		service.ColorService.Print(answer)
	}
}

func getModuleConfiguration(moduleName string) (*cfg.Config, []byte) {
	if moduleName == "" {
		fmt.Println("Need module name. Use 'info' for get module names")
		return nil, nil
	}

	moduleConfiguration, err := service.ConfigClient.GetConfigByModuleName(moduleName)
	if err != nil {
		if errorResponse, ok := err.(http.ErrorResponse); ok {
			if errorResponse.StatusCode == fasthttp.StatusNotFound {
				fmt.Printf("Module '%s' not found\n", moduleName)
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
		return nil, nil
	}

	if jsonObject, err := json.Marshal(moduleConfiguration.Data); err != nil {
		fmt.Println(err)
		return nil, nil
	} else {
		return moduleConfiguration, jsonObject
	}
}

func createUpdateConfig(stringToChange string, configuration *cfg.Config) {
	newData := make(map[string]interface{})
	if stringToChange != "" {
		if err := json.Unmarshal([]byte(stringToChange), &newData); err != nil {
			fmt.Println(err)
			return
		}
	}
	configuration.Data = newData
	if resp, err := service.ConfigClient.CreateUpdateConfig(*configuration); err != nil {
		fmt.Println(err)
	} else {
		printAnswer(resp.Data)
	}
}
