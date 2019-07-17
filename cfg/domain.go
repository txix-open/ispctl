package cfg

import (
	"time"
)

type getModuleByUuidAndNameRequest struct {
	Name string
	Uuid string
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
