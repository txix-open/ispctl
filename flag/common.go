package flag

import "github.com/codegangsta/cli"

var (
	Host  = cli.StringFlag{Name: GateHostName}
	Uuid  = cli.StringFlag{Name: InstanceUuidName}
	Color = cli.BoolFlag{Name: ColorName}
)
