package utils

import (
	"isp-ctl/cfg"
	"isp-ctl/service"
)

func CreateUpdateConfig(stringToChange string, configuration *cfg.Config) {
	if answer, err := service.Config.CreateUpdateConfig(stringToChange, configuration); err != nil {
		PrintError(err)
	} else if answer != nil {
		PrintAnswer(answer)
	}
}
