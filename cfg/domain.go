package cfg

import (
	"time"

	"github.com/integration-system/isp-kit/cluster"
	"github.com/integration-system/isp-kit/rc/schema"
)

type getModuleByUuidAndNameRequest struct {
	ModuleName string
	Uuid       string
}

type getSchemaByModuleIdRequest struct {
	ModuleId string
}

type Config struct {
	Id            string
	Name          string
	CommonConfigs []string
	Description   string
	ModuleId      string
	Version       int32
	Active        bool
	Data          map[string]interface{}
	Unsafe        bool
}

type CommonConfig struct {
	Id          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Data        map[string]interface{}
}

type Connection struct {
	LibVersion    string
	Version       string
	Address       cluster.AddressConfiguration
	Endpoints     []cluster.EndpointDescriptor `json:",omitempty"`
	EstablishedAt time.Time
}

type ModuleInfo struct {
	Id                 string
	Name               string
	Active             bool
	CreatedAt          time.Time
	LastConnectedAt    time.Time
	LastDisconnectedAt time.Time
	Configs            []Config       `json:",omitempty"`
	ConfigSchema       *schema.Schema `json:",omitempty"`
	Status             []Connection   `json:",omitempty"`
}

type ConfigSchema struct {
	Id        string
	Version   string
	ModuleId  string
	Schema    schema.Schema
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Deleted struct {
	Deleted bool
	Links   Links
}

type Links map[string][]string

type identityRequest struct {
	Id string
}

type CompileConfigs struct {
	Data                map[string]interface{}
	CommonConfigsIdList []string
}
