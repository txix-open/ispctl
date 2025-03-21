package utils

import (
	"ispctl/model"
)

type UpdateConfigService interface {
	CreateUpdateConfig(stringToChange string, configuration *model.Config) (map[string]any, error)
}

func CreateUpdateConfig(stringToChange string, configuration *model.Config, service UpdateConfigService) error {
	answer, err := service.CreateUpdateConfig(stringToChange, configuration)
	if err != nil {
		return err
	} else if answer != nil {
		PrintAnswer(answer)
	}
	return nil
}
