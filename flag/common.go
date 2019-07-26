package flag

import "github.com/codegangsta/cli"

var (
	Host   = cli.StringFlag{Name: GateHostName, Usage: "overrides gateHost"}
	Uuid   = cli.StringFlag{Name: InstanceUuidName, Usage: "overrides instanceUuid"}
	Color  = cli.BoolFlag{Name: ColorName, Usage: "colorize the json for outputing to screen"}
	Unsafe = cli.BoolFlag{Name: UnsafeName, Usage: "disable checking schema for set configuration"}
)
