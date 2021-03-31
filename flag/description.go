package flag

import "fmt"

const (
	hostName  = "gateHost"
	hostUsage = "overrides gateHost; config-service address, default '127.0.0.1:9002'"

	uuidName  = "instanceUuid"
	uuidUsage = "overrides instanceUuid"

	colorName  = "color"
	colorUsage = "colorize the json for outputing to screen"

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
