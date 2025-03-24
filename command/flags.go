package command

import "fmt"

const (
	HostFlagName  = "gateHost"
	HostFlagUsage = "isp-config-service address"

	UnsafeFlagName  = "unsafe"
	UnsafeFlagUsage = "disable checking schema for set configuration"

	OutPrintFlagName  = "o"
	OutPrintJsonValue = "json"
	OutPrintHtmlValue = "html"

	SecretFlagName          = "secret"
	secretVariableFlagUsage = "set variable as secret"
)

var (
	outPrintStatusUsage = fmt.Sprintf("out print for status; %s", OutPrintJsonValue)
	outPrintSchemaUsage = fmt.Sprintf("out print for schema; %s or %s", OutPrintJsonValue, OutPrintHtmlValue)
)
