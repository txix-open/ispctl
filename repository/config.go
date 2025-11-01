package repository

import (
	"context"
	"fmt"
	"ispctl/model"
	"strings"
	"time"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"

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

func (c Config) GetAvailableConfigs() ([]model.ModuleInfo, error) {
	response := make([]model.ModuleInfo, 0)
	err := c.cli.Invoke(getAvailableConfigs).
		JsonResponseBody(&response).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getAvailableConfigs)
	}
	return response, nil
}

func (c Config) GetConfigByModuleName(name string) (*model.Config, error) {
	request := &model.GetModuleByUuidAndNameRequest{ModuleName: name}
	response := new(model.Config)
	err := c.cli.Invoke(getConfigByModuleName).
		JsonRequestBody(request).
		JsonResponseBody(&response).
		Timeout(time.Second).
		Do(c.baseCtx)
	switch {
	case status.Code(err) == codes.NotFound:
		return nil, model.ErrModuleNotFound
	case err != nil:
		return nil, errors.WithMessagef(err, "call %s", getConfigByModuleName)
	default:
		return response, nil
	}
}

func (c Config) CreateUpdateConfig(request model.Config) (*model.Config, error) {
	response := new(model.Config)
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

func (c Config) GetSchemaByModuleId(moduleId string) (*model.ConfigSchema, error) {
	request := model.GetSchemaByModuleIdRequest{ModuleId: moduleId}
	response := new(model.ConfigSchema)
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

func (c Config) GetAllVariables() ([]model.Variable, error) {
	response := make([]model.Variable, 0)
	err := c.cli.Invoke(getAllVariables).
		JsonResponseBody(&response).
		Timeout(time.Second).
		Do(c.baseCtx)
	if err != nil {
		return nil, errors.WithMessagef(err, "call %s", getAllVariables)
	}
	return response, nil
}

func (c Config) GetVariableByName(variableName string) (*model.Variable, error) {
	request := model.VariableByNameRequest{Name: variableName}
	response := new(model.Variable)
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

func (c Config) UpsertVariables(request []model.UpsertVariableRequest) error {
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
	request := model.VariableByNameRequest{Name: variableName}
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
	switch apiError.ErrorCode {
	case model.ErrCodeVariableNotFound:
		return model.ErrVariableNotFound
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

	apiErr := apierrors.FromError(err)
	if apiErr != nil {
		for field, desc := range apiErr.Details {
			stringBuilder.WriteString(fmt.Sprintf("%s | %s\n", field, desc))
		}
	} else {
		for _, detail := range errorStatus.Details() {
			switch d := detail.(type) {
			case *epb.BadRequest:
				for _, violation := range d.FieldViolations {
					stringBuilder.WriteString(fmt.Sprintf("%s | %s\n", violation.Field, violation.Description))
				}
			default:
				stringBuilder.WriteString(fmt.Sprint(d))
			}
		}
	}

	return errors.New(stringBuilder.String())
}
