package service

import (
	"encoding/json"
	"fmt"
	"os"

	"ispctl/cfg"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/grpc/apierrors"
	"github.com/txix-open/isp-kit/rc/schema"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var Config configService

type configService struct {
	UnsafeEnable bool
}

func (c configService) ReceiveConfiguration(host string) error {
	return ConfigClient.ReceiveConfig(host)
}

func (c configService) GetAvailableConfigs() ([]cfg.ModuleInfo, error) {
	return ConfigClient.GetAvailableConfigs()
}

func (c configService) GetConfigurationByModuleName(moduleName string) (*cfg.Config, error) {
	if moduleName == "" {
		return nil, errors.New("need module name")
	}

	if moduleConfiguration, err := ConfigClient.GetConfigByModuleName(moduleName); err != nil {
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
	if configSchema, err := ConfigClient.GetSchemaByModuleId(moduleId); err != nil {
		return nil, err
	} else {
		return configSchema.Schema, nil
	}
}

func (c configService) CreateUpdateConfig(stringToChange string, configuration *cfg.Config) (map[string]interface{}, error) {
	newData := make(map[string]interface{})
	if stringToChange != "" {
		if err := json.Unmarshal([]byte(stringToChange), &newData); err != nil {
			return nil, err
		}
	}

	configuration.Data = newData
	return c.CreateUpdateConfigV2(configuration)
}

func (c configService) CreateUpdateConfigV2(configuration *cfg.Config) (map[string]interface{}, error) {
	configuration.Unsafe = c.UnsafeEnable
	if resp, err := ConfigClient.CreateUpdateConfig(*configuration); err != nil {
		if errorStatus, ok := status.FromError(err); ok {
			switch errorStatus.Code() {
			case codes.InvalidArgument:
				fmt.Print("doesn't match with schema:\n")
				defer func() {
					os.Exit(1)
				}()
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

func (configService) GetAllVariables() ([]cfg.Variable, error) {
	return ConfigClient.GetAllVariables()
}

func (configService) GetVariableByName(variableName string) (*cfg.Variable, error) {
	variable, err := ConfigClient.GetVariableByName(variableName)
	if err == nil {
		return variable, nil
	}

	apiError := apierrors.FromError(err)
	if apiError != nil && apiError.ErrorCode == cfg.ErrCodeVariableNotFound {
		return nil, errors.Errorf("variable [%s] not found", variableName)
	}
	return nil, err
}

func (configService) DeleteVariable(variableName string) error {
	err := ConfigClient.DeleteVariable(variableName)
	apiError := apierrors.FromError(err)
	if apiError != nil && apiError.ErrorCode == cfg.ErrCodeVariableNotFound {
		return errors.Errorf("variable [%s] not found", variableName)
	}
	return err
}

func (configService) UpsertVariables(vars []cfg.UpsertVariableRequest) error {
	return ConfigClient.UpsertVariables(vars)
}
