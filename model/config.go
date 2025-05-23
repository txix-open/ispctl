package model

import (
	"time"

	"github.com/txix-open/isp-kit/cluster"
	"github.com/txix-open/isp-kit/rc/schema"
)

type GetModuleByUuidAndNameRequest struct {
	ModuleName string
	Uuid       string
}

type GetSchemaByModuleIdRequest struct {
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
	Data          map[string]any
	Unsafe        bool
}

type CommonConfig struct {
	Id          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Data        map[string]any
}

type Connection struct {
	LibVersion    string
	Version       string
	Address       cluster.AddressConfiguration
	Endpoints     []cluster.EndpointDescriptor `json:",omitempty"`
	EstablishedAt time.Time
}

type ModuleInfo struct {
	Id           string
	Name         string
	Active       bool
	CreatedAt    time.Time
	ConfigSchema *schema.Schema `json:",omitempty"`
	Status       []Connection   `json:",omitempty"`
}

type ConfigSchema struct {
	Id        string
	Version   string
	ModuleId  string
	Schema    schema.Schema
	CreatedAt time.Time
	UpdatedAt time.Time
}
