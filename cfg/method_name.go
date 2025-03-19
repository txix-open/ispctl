package cfg

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
