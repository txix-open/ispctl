package utils

import (
	"ispctl/cfg"
	"ispctl/service"
)

func CreateUpdateCommonConfig(stringToChange string, configuration cfg.CommonConfig) error {
	if answer, err := service.Config.CreateUpdateCommonConfig(stringToChange, configuration); err != nil {
		return err
	} else if answer != nil {
		PrintAnswer(answer)
	}
	return nil
}
