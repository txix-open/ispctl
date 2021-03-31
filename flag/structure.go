package flag

import (
	"strings"

	"github.com/urfave/cli/v2"
	"isp-ctl/service"
)

var (
	//---global---
	Host   = &cli.StringFlag{Name: hostName, Usage: hostUsage, Aliases: []string{"g", "configAddr"}}
	Uuid   = &cli.StringFlag{Name: uuidName, Usage: uuidUsage, Aliases: []string{"u"}}
	Color  = &cli.BoolFlag{Name: colorName, Usage: colorUsage, Aliases: []string{"c"}}
	Unsafe = &cli.BoolFlag{Name: unsafeName, Usage: unsafeUsage}
	//---local---
	OutPrintStatus   = &cli.StringFlag{Name: outPrintName, Usage: outPrintStatusUsage}
	OutPrintSchema   = &cli.StringFlag{Name: outPrintName, Usage: outPrintSchemaUsage, Value: OutPrintJsonValue}
	WithCommonConfig = &cli.BoolFlag{Name: withCommonConfigName, Usage: withCommonConfigUsage}
)

func CheckGlobal(c *cli.Context) error {
	var (
		uuid, host string
		color      bool
	)
	if c.Bool(Color.Name) != false {
		color = true
	}
	service.ColorService.Enable = color
	service.Config.UnsafeEnable = c.Bool(Unsafe.Name)

	if c.String(Host.Name) != "" {
		host = c.String(Host.Name)
	}

	if c.String(Uuid.Name) != "" {
		uuid = c.String(Uuid.Name)
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
