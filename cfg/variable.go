package cfg

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
	Description string `validate:"required"`
	Type        string `validate:"required,oneof=SECRET TEXT"`
	Value       string `validate:"required"`
}

type VariableByNameRequest struct {
	Name string `validate:"required"`
}
