package utils

import (
	"ispctl/entity"
	"ispctl/service"
)

func CreateUpdateConfig(stringToChange string, configuration *entity.Config) error {
	answer, err := service.Config.CreateUpdateConfig(stringToChange, configuration)
	if err != nil {
		return err
	} else if answer != nil {
		PrintAnswer(answer)
	}
	return nil
}
