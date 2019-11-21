package service

import (
	"encoding/json"
	"fmt"
	"github.com/integration-system/isp-lib/config/schema"
	"github.com/pkg/errors"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
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

	configuration.Unsafe = c.UnsafeEnable
	configuration.Data = newData
	if resp, err := configClient.CreateUpdateConfig(*configuration); err != nil {
		if errorStatus, ok := status.FromError(err); ok {
			switch errorStatus.Code() {
			case codes.InvalidArgument:
				fmt.Print("doesn't match with schema:\n")
				for _, value := range errorStatus.Details() {
					if a, ok := value.(*epb.BadRequest); ok {
						for _, value := range a.FieldViolations {
							fmt.Printf("%s | %s\n", value.Field, value.Description)
						}
					} else {
						fmt.Println(value)
					}
					return nil, nil
				}
			}
		}
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

func (c configService) GetCommonConfigByName(configName string) (cfg.CommonConfig, error) {
	if commonConfigs, err := c.GetMapNameCommonConfig(); err != nil {
		return cfg.CommonConfig{}, err
	} else if config, ok := commonConfigs[configName]; !ok {
		return cfg.CommonConfig{}, errors.Errorf("common config [%s] not found", configName)
	} else {
		return config, nil
	}
}

func (c configService) GetMapNameCommonConfig() (map[string]cfg.CommonConfig, error) {
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

func (c configService) UnlinkCommonConfigFromModule(configName, moduleName string) error {
	return c.linked(configName, moduleName,
		func(configId string, configIdentities []string) (bool, []string) {
			needUnlink := false
			for i, id := range configIdentities {
				if configId == id {
					copy(configIdentities[i:], configIdentities[i+1:])
					configIdentities[len(configIdentities)-1] = ""
					configIdentities = configIdentities[:len(configIdentities)-1]
					needUnlink = true
					break
				}
			}
			return needUnlink, configIdentities
		})
}

func (c configService) LinkCommonConfigToModule(configName, moduleName string) error {
	return c.linked(configName, moduleName,
		func(configId string, configIdentities []string) (bool, []string) {
			needLink := true
			for _, id := range configIdentities {
				if configId == id {
					needLink = false
					break
				}
			}
			if needLink {
				configIdentities = append(configIdentities, configId)
			}
			return needLink, configIdentities
		})
}

func (c configService) linked(configName, moduleName string,
	linkHandler func(configId string, configIdentities []string) (bool, []string)) error {

	if configName == "" {
		return errors.New("empty config name")
	}

	if moduleName == "" {
		return errors.New("empty module name")
	}

	commonConfigs, err := c.GetMapNameCommonConfig()
	if err != nil {
		return err
	}

	config, ok := commonConfigs[configName]
	if !ok {
		return errors.Errorf("common config [%s] not found", configName)
	}

	moduleConfig, err := c.GetConfigurationByModuleName(moduleName)
	if err != nil {
		return err
	}

	configIdConfigNameMap := make(map[string]string)
	for name, value := range commonConfigs {
		configIdConfigNameMap[value.Id] = name
	}

	needUpdate, newCommonConfigs := linkHandler(config.Id, moduleConfig.CommonConfigs)

	if needUpdate {
		moduleConfig.CommonConfigs = newCommonConfigs
		if linked, err := c.UpdateConfigAndGetLinkCommon(*moduleConfig); err != nil {
			return err
		} else {
			fmt.Printf("module [%s] have next link:\n", moduleName)
			for _, name := range linked {
				fmt.Printf("[%s] ", configIdConfigNameMap[name])
			}
			fmt.Printf("\n")
		}
	} else {
		fmt.Printf("module [%s] have next link:\n", moduleName)
		for _, name := range moduleConfig.CommonConfigs {
			fmt.Printf("[%s] ", configIdConfigNameMap[name])
		}
		fmt.Printf("\n")
	}
	return nil
}
