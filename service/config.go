package service

import (
	"encoding/json"
	"fmt"
	"os"

	"ispctl/entity"
	"ispctl/repository"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/grpc/apierrors"
	"github.com/txix-open/isp-kit/rc/schema"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var Config = &ConfigService{}

type ConfigService struct {
	UnsafeEnable bool
	configRepo   repository.Config
}

func (c *ConfigService) ReceiveConfiguration(host string) error {
	cli, err := repository.NewClientFromHost(host)
	if err != nil {
		return err
	}
	c.configRepo = repository.NewConfig(cli)
	return nil
}

func (c *ConfigService) GetAvailableConfigs() ([]entity.ModuleInfo, error) {
	return c.configRepo.GetAvailableConfigs()
}

func (c *ConfigService) GetConfigurationByModuleName(moduleName string) (*entity.Config, error) {
	if moduleName == "" {
		return nil, errors.New("need module name")
	}

	moduleConfiguration, err := c.configRepo.GetConfigByModuleName(moduleName)
	if err == nil {
		return moduleConfiguration, nil
	}
	errorStatus, ok := status.FromError(err)
	if ok && errorStatus.Code() == codes.NotFound {
		return nil, errors.Errorf("module [%s] not found", moduleName)
	}
	return nil, err
}

func (c *ConfigService) GetSchemaByModuleId(moduleId string) (schema.Schema, error) {
	configSchema, err := c.configRepo.GetSchemaByModuleId(moduleId)
	if err != nil {
		return nil, err
	}
	return configSchema.Schema, nil
}

func (c *ConfigService) CreateUpdateConfig(stringToChange string, configuration *entity.Config) (map[string]any, error) {
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

func (c *ConfigService) CreateUpdateConfigV2(configuration *entity.Config) (map[string]any, error) {
	configuration.Unsafe = c.UnsafeEnable
	resp, err := c.configRepo.CreateUpdateConfig(*configuration)
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
			switch errorDetails := value.(type) {
			case *epb.BadRequest:
				for _, value := range errorDetails.FieldViolations {
					fmt.Printf("%s | %s\n", value.Field, value.Description)
				}
			default:
				fmt.Println(value)
			}
		}
		return nil, nil
	}
	return nil, err
}

func (c *ConfigService) GetAllVariables() ([]entity.Variable, error) {
	return c.configRepo.GetAllVariables()
}

func (c *ConfigService) GetVariableByName(variableName string) (*entity.Variable, error) {
	if variableName == "" {
		return nil, errors.New("need variable name")
	}

	variable, err := c.configRepo.GetVariableByName(variableName)
	if err == nil {
		return variable, nil
	}

	apiError := apierrors.FromError(err)
	if apiError != nil && apiError.ErrorCode == entity.ErrCodeVariableNotFound {
		return nil, errors.Errorf("variable [%s] not found", variableName)
	}
	return nil, err
}

func (c *ConfigService) DeleteVariable(variableName string) error {
	if variableName == "" {
		return errors.New("need variable name")
	}

	err := c.configRepo.DeleteVariable(variableName)
	apiError := apierrors.FromError(err)
	if apiError != nil && apiError.ErrorCode == entity.ErrCodeVariableNotFound {
		return errors.Errorf("variable [%s] not found", variableName)
	}
	return err
}

func (c *ConfigService) UpsertVariables(vars []entity.UpsertVariableRequest) error {
	return c.configRepo.UpsertVariables(vars)
}
