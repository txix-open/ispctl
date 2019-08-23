package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/integration-system/isp-lib/config/schema"
	"github.com/integration-system/isp-lib/http"
	"github.com/valyala/fasthttp"
	"github.com/xeipuuv/gojsonschema"
	"isp-ctl/cfg"
)

var (
	configClient = cfg.NewConfigClient(http.NewJsonRestClient())
	Config       configService
)

type configService struct {
	UnsafeEnable bool
}

func (configService) ReceiveConfiguration(host, uuid string) {
	configClient.ReceiveConfig(host, uuid)
}

func (configService) GetAvailableConfigs() ([]cfg.ModuleInfo, error) {
	return configClient.GetAvailableConfigs()
}

func (configService) GetConfigurationAndJsonByModuleName(moduleName string) (*cfg.Config, []byte, error) {
	if moduleName == "" {
		return nil, nil, errors.New("need module name")
	}
	if moduleConfiguration, err := configClient.GetConfigByModuleName(moduleName); err != nil {
		if errorResponse, ok := err.(http.ErrorResponse); ok && errorResponse.StatusCode == fasthttp.StatusNotFound {
			return nil, nil, errors.New(fmt.Sprintf("module '%s' not found", moduleName))
		}
		return nil, nil, err
	} else if jsonObject, err := json.Marshal(moduleConfiguration.Data); err != nil {
		return nil, nil, err
	} else {
		return moduleConfiguration, jsonObject, nil
	}
}

func (configService) GetSchemaByModuleId(moduleId int32) (schema.Schema, error) {
	if configSchema, err := configClient.GetSchemaByModuleId(moduleId); err != nil {
		return nil, err
	} else {
		dereferenceSchema := schema.DereferenceSchema(configSchema.Schema)
		return dereferenceSchema, nil
	}
}

func (c configService) CreateUpdateConfig(stringToChange string, configuration *cfg.Config) (map[string]interface{}, error) {
	newData := make(map[string]interface{})
	if stringToChange != "" {
		if err := json.Unmarshal([]byte(stringToChange), &newData); err != nil {
			return nil, err
		}
	}

	if ok, err := c.checkSchema(configuration.ModuleId, newData); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	configuration.Data = newData
	if resp, err := configClient.CreateUpdateConfig(*configuration); err != nil {
		return nil, err
	} else {
		return resp.Data, nil
	}
}

func (c configService) checkSchema(moduleId int32, object map[string]interface{}) (bool, error) {
	if c.UnsafeEnable {
		return true, nil
	}

	if resp, err := configClient.GetSchemaByModuleId(moduleId); err != nil {
		return false, err
	} else {
		schemaLoader := gojsonschema.NewGoLoader(resp.Schema)
		documentLoader := gojsonschema.NewGoLoader(object)

		if result, err := gojsonschema.Validate(schemaLoader, documentLoader); err != nil {
			return false, err
		} else if result.Valid() {
			return true, nil
		} else {
			for _, desc := range result.Errors() {
				fmt.Printf("- %s\n", desc)
			}
			return false, nil
		}
	}
}
