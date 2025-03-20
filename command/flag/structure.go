package flag

import (
	"strings"

	"ispctl/service"

	"github.com/urfave/cli/v2"
)

var (
	//---global---
	Host = &cli.StringFlag{
		Name:     hostName,
		Usage:    hostUsage,
		Value:    "127.0.0.1:9002",
		Required: false,
		Aliases:  []string{"g", "configAddr"},
	}
	Unsafe = &cli.BoolFlag{Name: unsafeName, Usage: unsafeUsage}
	//---local---
	OutPrintStatus = &cli.StringFlag{Name: outPrintName, Usage: outPrintStatusUsage}
	OutPrintSchema = &cli.StringFlag{Name: outPrintName, Usage: outPrintSchemaUsage, Value: OutPrintJsonValue}

	SetVariableSecretType = &cli.BoolFlag{Name: secretVariableName, Usage: secretVariableUsage}
)

func ApplyGlobalFlags(c *cli.Context) error {
	var host string
	service.Config.UnsafeEnable = c.Bool(Unsafe.Name)

	host = c.String(Host.Name)
	host = strings.Replace(host, "'", "", -1)

	return service.Config.ReceiveConfiguration(host)
}
