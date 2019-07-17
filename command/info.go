package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/olekukonko/tablewriter"
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

func (statusCommand) action(c *cli.Context) {
	if err := checkFlags(c); err != nil {
		fmt.Println(err)
		return
	}
	if arrayOfModules, err := service.ConfigClient.GetAvailableConfigs(); err != nil {
		fmt.Println(err)
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
		table.SetHeader([]string{tableHeaderModuleName, tableHeaderStatus, tableHeaderAddresses})
		for _, module := range arrayOfModules {
			arrayOfAddress := ""
			connection := tableNotConnectedStatus
			for _, value := range module.Status {
				if arrayOfAddress == "" {
					arrayOfAddress = fmt.Sprintf("%s", value.Address.IP)
				} else {
					arrayOfAddress = fmt.Sprintf("%s, %s", arrayOfAddress, value.Address)
				}
			}
			if arrayOfAddress != "" {
				connection = tableConnectedStatus
			}
			table.Append([]string{module.Name, connection, arrayOfAddress})
		}
		table.Render()
	}
}
