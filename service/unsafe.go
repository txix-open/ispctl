package service

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
)

var UnsafeService unsafeService

type unsafeService struct {
	Enable bool
}

func (s unsafeService) CheckConfigurationSchema(moduleId int32, object map[string]interface{}) (bool, error) {
	if s.Enable {
		return true, nil
	}

	if resp, err := ConfigClient.GetSchemaByModuleId(moduleId); err != nil {
		return false, err
	} else {
		schemaLoader := gojsonschema.NewGoLoader(resp.Schema)
		documentLoader := gojsonschema.NewGoLoader(object)

		if result, err := gojsonschema.Validate(schemaLoader, documentLoader); err != nil {
			return false, err
		} else if result.Valid() {
			return true, nil
		} else {
			for _, desc := range result.Errors() {
				fmt.Printf("- %s\n", desc)
			}
			return false, nil
		}
	}
}
