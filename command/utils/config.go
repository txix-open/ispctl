package utils

import (
	"ispctl/model"
	"ispctl/service"
)

func CreateUpdateConfig(stringToChange string, configuration *model.Config) error {
	answer, err := service.Config.CreateUpdateConfig(stringToChange, configuration)
	if err != nil {
		return err
	} else if answer != nil {
		PrintAnswer(answer)
	}
	return nil
}
