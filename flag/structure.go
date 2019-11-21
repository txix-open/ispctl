package flag

import (
	"github.com/codegangsta/cli"
	"isp-ctl/service"
	"strings"
)

var (
	//---global---
	Host   = cli.StringFlag{Name: hostName, Usage: hostUsage}
	Uuid   = cli.StringFlag{Name: uuidName, Usage: uuidUsage}
	Color  = cli.BoolFlag{Name: colorName, Usage: colorUsage}
	Unsafe = cli.BoolFlag{Name: unsafeName, Usage: unsafeUsage}
	//---local---
	OutPrint         = cli.StringFlag{Name: outPrintName, Usage: outPrintUsage, Value: OutPrintJsonValue}
	WithCommonConfig = cli.BoolFlag{Name: withCommonConfigName, Usage: withCommonConfigUsage}
)

func CheckGlobal(c *cli.Context) error {
	var (
		uuid, host string
		color      bool
	)
	colorFlag := strings.Split(Color.Name, ", ")
	if c.GlobalBool(colorFlag[0]) != false {
		color = true
	} else {
		color = c.GlobalBool(colorFlag[1])
	}
	service.ColorService.Enable = color
	service.Config.UnsafeEnable = c.GlobalBool(Unsafe.Name)

	hostFlag := strings.Split(Host.Name, ", ")
	if c.GlobalString(hostFlag[0]) != "" {
		host = c.GlobalString(hostFlag[0])
	} else {
		host = c.GlobalString(hostFlag[1])
	}

	uuidFlag := strings.Split(Uuid.Name, ", ")
	if c.GlobalString(uuidFlag[0]) != "" {
		uuid = c.GlobalString(uuidFlag[0])
	} else {
		uuid = c.GlobalString(uuidFlag[1])
	}

	uuid = strings.Replace(uuid, "'", "", -1)
	host = strings.Replace(host, "'", "", -1)
	if uuid == "" || host == "" {
		if configuration, err := service.YamlService.Parse(); err != nil {
			return err
		} else {
			if uuid == "" {
				uuid = configuration.InstanceUuid
			}
			if host == "" {
				host = configuration.GateHost
			}
		}
	}
	return service.Config.ReceiveConfiguration(host, uuid)
}
