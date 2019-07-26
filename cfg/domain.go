package cfg

import (
	"github.com/integration-system/isp-lib/config/schema"
	"github.com/integration-system/isp-lib/structure"
	"time"
)

type getModuleByUuidAndNameRequest struct {
	Name string
	Uuid string
}

type getSchemaByModuleIdRequest struct {
	ModuleId int32
}

type Config struct {
	Id          int64
	Name        string
	Description string
	ModuleId    int32
	Version     int32
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Data        map[string]interface{}
}

type Connection struct {
	LibVersion    string
	Version       string
	Address       structure.AddressConfiguration
	Endpoints     []structure.EndpointConfig `json:",omitempty"`
	EstablishedAt time.Time
}

type ModuleInfo struct {
	Id                 int32
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
	Id        int32
	Version   string
	ModuleId  int32
	Schema    schema.Schema
	CreatedAt time.Time
	UpdatedAt time.Time
}
