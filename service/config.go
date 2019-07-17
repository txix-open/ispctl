package service

import (
	"github.com/integration-system/isp-lib/http"
	"isp-ctl/cfg"
)

var ConfigClient = cfg.NewConfigClient(http.NewJsonRestClient())

var ConfigService configService

type configService struct{}

func (configService) ReceiveConfiguration(host, uuid string) {
	ConfigClient.ReceiveConfig(host, uuid)
}
