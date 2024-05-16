package utils

import (
	"ispctl/cfg"
	"ispctl/service"
)

func CreateUpdateConfig(stringToChange string, configuration *cfg.Config) error {
	if answer, err := service.Config.CreateUpdateConfig(stringToChange, configuration); err != nil {
		return err
	} else if answer != nil {
		PrintAnswer(answer)
	}
	return nil
}
