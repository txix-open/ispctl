package model

import "github.com/pkg/errors"

var (
	ErrModuleNotFound   = errors.New("module not found")
	ErrVariableNotFound = errors.New("variable not found")
)

const (
	SecretVariableType = "SECRET"
	TextVariableType   = "TEXT"

	ErrCodeVariableNotFound = 2006
)

type Variable struct {
	Name              string
	Description       string
	Type              string
	Value             string
	ContainsInConfigs []LinkedConfig
}

type LinkedConfig struct {
	Id       string
	ModuleId string
	Name     string
}

type UpsertVariableRequest struct {
	Name        string `validate:"required"`
	Description string
	Type        string `validate:"required,oneof=SECRET TEXT"`
	Value       string `validate:"required"`
}

type VariableByNameRequest struct {
	Name string `validate:"required"`
}
