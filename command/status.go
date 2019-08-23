package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/olekukonko/tablewriter"
	"isp-ctl/flag"
	"isp-ctl/service"
	"os"
)

const (
	tableHeaderModuleName = "MODULE_NAME"
	tableHeaderStatus     = "STATUS"
	tableHeaderAddresses  = "ADDRESSES"

	tableNotConnectedStatus = "NOT_CONNECTED"
	tableConnectedStatus    = "CONNECTED"
)

func Status() cli.Command {
	return cli.Command{
		Name:   "status",
		Usage:  "get available configs",
		Action: status.action,
	}
}

var status statusCommand

type statusCommand struct{}

func (statusCommand) action(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		printError(err)
		return
	}
	if arrayOfModules, err := service.Config.GetAvailableConfigs(); err != nil {
		printError(err)
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
		table.SetHeader([]string{tableHeaderModuleName, tableHeaderStatus, tableHeaderAddresses})
		for _, module := range arrayOfModules {
			addresses := ""
			connection := tableNotConnectedStatus
			for _, value := range module.Status {
				if value.Address.IP == "" && value.Address.Port == "" {
					continue
				}
				if addresses == "" {
					addresses = fmt.Sprintf("%s:%s;", value.Address.IP, value.Address.Port)
				} else {
					addresses = fmt.Sprintf("%s\n%s:%s;", addresses, value.Address.IP, value.Address.Port)
				}
			}
			if addresses != "" {
				connection = tableConnectedStatus
			}
			table.Append([]string{module.Name, connection, addresses})
		}
		table.Render()
	}
}
