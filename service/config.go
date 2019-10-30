package service

import (
	"encoding/json"
	"fmt"
	"github.com/integration-system/isp-lib/config/schema"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"isp-ctl/cfg"
)

var Config configService

type configService struct {
	UnsafeEnable bool
}

func (c configService) ReceiveConfiguration(host, uuid string) error {
	return configClient.ReceiveConfig(host, uuid)
}

func (c configService) GetAvailableConfigs() ([]cfg.ModuleInfo, error) {
	return configClient.GetAvailableConfigs()
}

func (c configService) GetConfigurationByModuleName(moduleName string) (*cfg.Config, error) {
	if moduleName == "" {
		return nil, errors.New("need module name")
	}

	if moduleConfiguration, err := configClient.GetConfigByModuleName(moduleName); err != nil {
		if errorStatus, ok := status.FromError(err); ok {
			switch errorStatus.Code() {
			case codes.NotFound:
				return nil, errors.Errorf("module [%s] not found", moduleName)
			}
		}
		return nil, err
	} else {
		return moduleConfiguration, nil
	}
}

func (c configService) GetSchemaByModuleId(moduleId string) (schema.Schema, error) {
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

func (c configService) UpdateConfigAndGetLinkCommon(configuration cfg.Config) ([]string, error) {
	if resp, err := configClient.CreateUpdateConfig(configuration); err != nil {
		return nil, err
	} else {
		return resp.CommonConfigs, nil
	}
}

func (c configService) GetMapCommonConfigByName() (map[string]cfg.CommonConfig, error) {
	collection, err := configClient.GetCommonConfigs([]string{})
	if err != nil {
		return nil, err
	}
	response := make(map[string]cfg.CommonConfig)
	for _, config := range collection {
		response[config.Name] = config
	}
	return response, nil
}

func (c configService) CreateUpdateCommonConfig(stringToChange string, config cfg.CommonConfig) (map[string]interface{}, error) {
	newData := make(map[string]interface{})
	if stringToChange != "" {
		if err := json.Unmarshal([]byte(stringToChange), &newData); err != nil {
			return nil, err
		}
	}

	config.Data = newData
	if resp, err := configClient.CreateUpdateCommonConfig(config); err != nil {
		return nil, err
	} else {
		return resp.Data, nil
	}
}

func (c configService) DeleteCommonConfig(configId string) ([]string, bool, error) {
	deletedInfo, err := configClient.DeleteCommonConfig(configId)
	if err != nil {
		return nil, false, err
	}
	arrayOfModuleName := make([]string, 0)
	for moduleName := range deletedInfo.Links {
		arrayOfModuleName = append(arrayOfModuleName, moduleName)
	}
	return arrayOfModuleName, deletedInfo.Deleted, nil
}

func (c configService) GetLinksCommonConfig(configId string) ([]string, error) {
	link, err := configClient.GetLinksCommonConfig(configId)
	if err != nil {
		return nil, err
	}
	arrayOfModuleName := make([]string, 0)
	for moduleName := range *link {
		arrayOfModuleName = append(arrayOfModuleName, moduleName)
	}
	return arrayOfModuleName, nil
}

func (c configService) CompileDataWithCommonConfigs(
	data map[string]interface{}, commonConfigs []string) (map[string]interface{}, error) {

	req := cfg.CompileConfigs{
		Data:                data,
		CommonConfigsIdList: commonConfigs,
	}
	var err error
	data, err = configClient.CompileCommonConfigs(req)
	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func (c configService) checkSchema(moduleId string, object map[string]interface{}) (bool, error) {
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
