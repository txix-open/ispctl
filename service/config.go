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

	moduleConfiguration, err := ConfigClient.GetConfigByModuleName(moduleName)
	if err == nil {
		return moduleConfiguration, nil
	}
	errorStatus, ok := status.FromError(err)
	if ok && errorStatus.Code() == codes.NotFound {
		return nil, errors.Errorf("module [%s] not found", moduleName)
	}
	return nil, err
}

func (c configService) GetSchemaByModuleId(moduleId string) (schema.Schema, error) {
	configSchema, err := ConfigClient.GetSchemaByModuleId(moduleId)
	if err != nil {
		return nil, err
	}
	return configSchema.Schema, nil
}

func (c configService) CreateUpdateConfig(stringToChange string, configuration *cfg.Config) (map[string]any, error) {
	newData := make(map[string]any)
	if stringToChange != "" {
		err := json.Unmarshal([]byte(stringToChange), &newData)
		if err != nil {
			return nil, err
		}
	}

	configuration.Data = newData
	return c.CreateUpdateConfigV2(configuration)
}

func (c configService) CreateUpdateConfigV2(configuration *cfg.Config) (map[string]any, error) {
	configuration.Unsafe = c.UnsafeEnable
	resp, err := ConfigClient.CreateUpdateConfig(*configuration)
	if err == nil {
		return resp.Data, nil
	}
	errorStatus, ok := status.FromError(err)
	if ok && errorStatus.Code() == codes.InvalidArgument {
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
	return nil, err
}

func (configService) GetAllVariables() ([]cfg.Variable, error) {
	return ConfigClient.GetAllVariables()
}

func (configService) GetVariableByName(variableName string) (*cfg.Variable, error) {
	if variableName == "" {
		return nil, errors.New("need variable name")
	}

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
	if variableName == "" {
		return errors.New("need variable name")
	}

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
