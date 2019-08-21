package flag

import "github.com/codegangsta/cli"

var (
	//---global---
	Host   = cli.StringFlag{Name: hostName, Usage: hostUsage}
	Uuid   = cli.StringFlag{Name: uuidName, Usage: uuidUsage}
	Color  = cli.BoolFlag{Name: colorName, Usage: colorUsage}
	Unsafe = cli.BoolFlag{Name: unsafeName, Usage: unsafeUsage}
	//---local---
	OutPrint = cli.StringFlag{Name: outPrintName, Usage: outPrintUsage, Value: OutPrintJsonValue}
)
