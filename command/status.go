package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/olekukonko/tablewriter"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
	"os"
	"strings"
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
		utils.PrintError(err)
		return
	}
	if arrayOfModules, err := service.Config.GetAvailableConfigs(); err != nil {
		utils.PrintError(err)
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
		table.SetHeader([]string{tableHeaderModuleName, tableHeaderStatus, tableHeaderAddresses})
		for _, module := range arrayOfModules {
			addresses := ""
			connection := tableNotConnectedStatus
			instanceList := make([]string, 0, len(module.Status))
			for _, value := range module.Status {
				instanceList = append(instanceList, fmt.Sprintf("%s %s", value.Address.GetAddress(), value.Version))
			}
			addresses = strings.Join(instanceList, "\n")
			if addresses != "" {
				connection = tableConnectedStatus
			}
			table.Append([]string{module.Name, connection, addresses})
		}
		table.Render()
	}
}
