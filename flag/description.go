package flag

import "fmt"

const (
	hostName  = "gateHost, g"
	hostUsage = "overrides gateHost"

	uuidName  = "instanceUuid, u"
	uuidUsage = "overrides instanceUuid"

	colorName  = "color, c"
	colorUsage = "colorize the json for outputing to screen"

	unsafeName  = "unsafe"
	unsafeUsage = "disable checking schema for set configuration"

	outPrintName      = "o"
	OutPrintJsonValue = "json"
	OutPrintHtmlValue = "html"
)

var (
	outPrintUsage = fmt.Sprintf("out print for schema; %s or %s", OutPrintJsonValue, OutPrintHtmlValue)
)
