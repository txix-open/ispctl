package cfg

import (
	"github.com/integration-system/isp-lib/backend"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	instanceIdHeader = "x-instance-identity"
	callerId         = 1
)

func NewConfigClient() *configClient {
	return &configClient{
		cli:     &backend.InternalGrpcClient{},
		headers: make(map[string][]string),
	}
}

type configClient struct {
	cli     *backend.InternalGrpcClient
	headers metadata.MD

	address      string
	instanceUuid string
}

func (c *configClient) ReceiveConfig(address, instanceUuid string) error {
	var err error
	c.cli, err = backend.NewGrpcClient(address, grpc.WithInsecure())
	if err != nil {
		return err
	}

	c.address = address
	c.instanceUuid = instanceUuid
	c.headers[instanceIdHeader] = []string{instanceUuid}
	return nil
}

func (c *configClient) GetAvailableConfigs() ([]ModuleInfo, error) {
	response := make([]ModuleInfo, 0)
	err := c.cli.Invoke(getAvailableConfigs, callerId, nil, &response, backend.WithMetadata(c.headers))
	if err != nil {
		return nil, c.errorHandler(err)
	} else {
		return response, nil
	}
}

func (c *configClient) GetConfigByModuleName(name string) (*Config, error) {
	request := &getModuleByUuidAndNameRequest{ModuleName: name, Uuid: c.instanceUuid}
	response := new(Config)
	err := c.cli.Invoke(getConfigByModuleName, callerId, request, response, backend.WithMetadata(c.headers))
	if err != nil {
		return nil, c.errorHandler(err)
	} else {
		return response, nil
	}
}

func (c *configClient) CreateUpdateConfig(request Config) (*Config, error) {
	response := new(Config)
	err := c.cli.Invoke(createUpdateConfig, callerId, &request, response, backend.WithMetadata(c.headers))
	if err != nil {
		return nil, c.errorHandler(err)
	} else {
		return response, nil
	}
}

func (c *configClient) GetSchemaByModuleId(moduleId string) (*ConfigSchema, error) {
	request := &getSchemaByModuleIdRequest{ModuleId: moduleId}
	response := new(ConfigSchema)
	err := c.cli.Invoke(getSchemaByModuleId, callerId, &request, response, backend.WithMetadata(c.headers))
	if err != nil {
		return nil, c.errorHandler(err)
	} else {
		return response, nil
	}
}

func (c *configClient) GetCommonConfigs(req []string) ([]CommonConfig, error) {
	response := make([]CommonConfig, 0)
	err := c.cli.Invoke(getCommonConfigs, callerId, req, &response, backend.WithMetadata(c.headers))
	if err != nil {
		return nil, c.errorHandler(err)
	} else {
		return response, nil
	}
}

func (c *configClient) CreateUpdateCommonConfig(req CommonConfig) (*CommonConfig, error) {
	response := new(CommonConfig)
	err := c.cli.Invoke(createUpdateCommonConfig, callerId, req, response, backend.WithMetadata(c.headers))
	if err != nil {
		return nil, c.errorHandler(err)
	} else {
		return response, nil
	}
}

func (c *configClient) DeleteCommonConfig(id string) (*Deleted, error) {
	request := identityRequest{Id: id}
	response := new(Deleted)
	err := c.cli.Invoke(deleteCommonConfig, callerId, request, response, backend.WithMetadata(c.headers))
	if err != nil {
		return nil, c.errorHandler(err)
	} else {
		return response, nil
	}
}

func (c *configClient) GetLinksCommonConfig(id string) (*Links, error) {
	request := identityRequest{Id: id}
	response := new(Links)
	err := c.cli.Invoke(getLinksCommonConfig, callerId, request, response, backend.WithMetadata(c.headers))
	if err != nil {
		return nil, c.errorHandler(err)
	} else {
		return response, nil
	}
}

func (c *configClient) CompileCommonConfigs(request CompileConfigs) (map[string]interface{}, error) {
	response := make(map[string]interface{})
	err := c.cli.Invoke(compileCommonConfigs, callerId, request, &response, backend.WithMetadata(c.headers))
	if err != nil {
		return nil, c.errorHandler(err)
	} else {
		return response, nil
	}
}

func (c *configClient) errorHandler(err error) error {
	errStatus, ok := status.FromError(err)
	if !ok {
		return err
	}
	switch errStatus.Code() {
	case codes.Unavailable:
		return errors.Errorf("config service didn't connection on %s", c.address)
	default:
		return err
	}
}
