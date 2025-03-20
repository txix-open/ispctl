package repository

import (
	"context"
	"fmt"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"ispctl/entity"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/grpc/apierrors"
	"github.com/txix-open/isp-kit/grpc/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	getAvailableConfigs   = "config/module/get_modules_info"
	getConfigByModuleName = "config/config/get_active_config_by_module_name"
	createUpdateConfig    = "config/config/create_update_config"
	getSchemaByModuleId   = "config/schema/get_by_module_id"

	getAllVariables   = "config/variable/all"
	getVariableByName = "config/variable/get_by_name"
	upsertVariable    = "config/variable/upsert"
	deleteVariable    = "config/variable/delete"
)

type Config struct {
	cli     *client.Client
	baseCtx context.Context
}

func NewConfig(cli *client.Client) Config {
	return Config{
		cli:     cli,
		baseCtx: context.Background(),
	}
}

func (c Config) GetAvailableConfigs() ([]entity.ModuleInfo, error) {
	response := make([]entity.ModuleInfo, 0)
	err := c.cli.Invoke(getAvailableConfigs).
		JsonResponseBody(&response).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getAvailableConfigs)
	}
	return response, nil
}

func (c Config) GetConfigByModuleName(name string) (*entity.Config, error) {
	request := &entity.GetModuleByUuidAndNameRequest{ModuleName: name}
	response := new(entity.Config)
	err := c.cli.Invoke(getConfigByModuleName).
		JsonRequestBody(request).
		JsonResponseBody(&response).
		Timeout(time.Second).
		Do(c.baseCtx)
	switch {
	case status.Code(err) == codes.NotFound:
		return nil, entity.ErrModuleNotFound
	case err != nil:
		return nil, errors.WithMessagef(err, "call %s", getConfigByModuleName)
	default:
		return response, nil
	}
}

func (c Config) CreateUpdateConfig(request entity.Config) (*entity.Config, error) {
	response := new(entity.Config)
	err := c.cli.Invoke(createUpdateConfig).
		JsonRequestBody(request).
		JsonResponseBody(&response).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return nil, c.handleCreateUpdateConfigError(err)
	}
	return response, nil
}

func (c Config) GetSchemaByModuleId(moduleId string) (*entity.ConfigSchema, error) {
	request := entity.GetSchemaByModuleIdRequest{ModuleId: moduleId}
	response := new(entity.ConfigSchema)
	err := c.cli.Invoke(getSchemaByModuleId).
		JsonRequestBody(request).
		JsonResponseBody(&response).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getSchemaByModuleId)
	}
	return response, nil
}

func (c Config) GetAllVariables() ([]entity.Variable, error) {
	response := make([]entity.Variable, 0)
	err := c.cli.Invoke(getAllVariables).
		JsonResponseBody(&response).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getAllVariables)
	}
	return response, nil
}

func (c Config) GetVariableByName(variableName string) (*entity.Variable, error) {
	request := entity.VariableByNameRequest{Name: variableName}
	response := new(entity.Variable)
	err := c.cli.Invoke(getVariableByName).
		JsonRequestBody(request).
		JsonResponseBody(&response).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return nil, c.handleVariableError(err, getVariableByName)
	}
	return response, nil
}

func (c Config) UpsertVariables(request []entity.UpsertVariableRequest) error {
	err := c.cli.Invoke(upsertVariable).
		JsonRequestBody(request).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return errors.WithMessagef(err, "call %s", upsertVariable)
	}
	return nil
}

func (c Config) DeleteVariable(variableName string) error {
	request := entity.VariableByNameRequest{Name: variableName}
	err := c.cli.Invoke(deleteVariable).
		JsonRequestBody(request).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return c.handleVariableError(err, deleteVariable)
	}
	return nil
}

func (c Config) handleVariableError(err error, endpoint string) error {
	apiError := apierrors.FromError(err)
	if apiError == nil {
		return errors.WithMessagef(err, "call %s", endpoint)
	}
	switch {
	case apiError.ErrorCode == entity.ErrCodeVariableNotFound:
		return entity.ErrVariableNotFound
	default:
		return errors.WithMessagef(err, "call %s", endpoint)
	}
}

func (c Config) handleCreateUpdateConfigError(err error) error {
	errorStatus, ok := status.FromError(err)
	if !ok || errorStatus.Code() != codes.InvalidArgument {
		return errors.WithMessagef(err, "call %s", createUpdateConfig)
	}
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString("doesn't match with schema:\n")
	for _, value := range errorStatus.Details() {
		switch errorDetails := value.(type) {
		case *epb.BadRequest:
			for _, value := range errorDetails.FieldViolations {
				stringBuilder.WriteString(fmt.Sprintf("%s | %s\n", value.Field, value.Description))
			}
		default:
			stringBuilder.WriteString(fmt.Sprint(value))
		}
	}
	return errors.New(stringBuilder.String())
}
