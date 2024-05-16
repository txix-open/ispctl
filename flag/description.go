package flag

import "fmt"

const (
	hostName  = "gateHost"
	hostUsage = "isp-config-service address"

	unsafeName  = "unsafe"
	unsafeUsage = "disable checking schema for set configuration"

	outPrintName      = "o"
	OutPrintJsonValue = "json"
	OutPrintHtmlValue = "html"

	withCommonConfigName  = "full"
	withCommonConfigUsage = "compile data with common configs"
)

var (
	outPrintStatusUsage = fmt.Sprintf("out print for status; %s", OutPrintJsonValue)
	outPrintSchemaUsage = fmt.Sprintf("out print for schema; %s or %s", OutPrintJsonValue, OutPrintHtmlValue)
)
