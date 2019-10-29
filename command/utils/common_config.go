package utils

import (
	"isp-ctl/cfg"
	"isp-ctl/service"
)

func CreateUpdateCommonConfig(stringToChange string, configuration cfg.CommonConfig) {
	if answer, err := service.Config.CreateUpdateCommonConfig(stringToChange, configuration); err != nil {
		PrintError(err)
	} else if answer != nil {
		PrintAnswer(answer)
	}
}
