package command

import (
	"fmt"
	"os"
	"strings"

	"ispctl/command/flag"
	"ispctl/command/utils"
	"ispctl/service"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

const (
	tableHeaderModuleName = "MODULE_NAME"
	tableHeaderStatus     = "STATUS"
	tableHeaderAddresses  = "ADDRESSES"

	tableNotConnectedStatus = "NOT_CONNECTED"
	tableConnectedStatus    = "CONNECTED"
)

func Status() *cli.Command {
	return &cli.Command{
		Name:   "status",
		Usage:  "get available configs",
		Before: flag.ApplyGlobalFlags,
		Action: status.action,
		Flags: []cli.Flag{
			flag.OutPrintStatus,
		},
	}
}

var status statusCommand

type statusCommand struct{}

func (statusCommand) action(ctx *cli.Context) error {
	if arrayOfModules, err := service.Config.GetAvailableConfigs(); err != nil {
		return err
	} else {
		switch ctx.String(flag.OutPrintStatus.Name) {
		case flag.OutPrintJsonValue:
			utils.PrintAnswer(arrayOfModules)
		default:
			table := tablewriter.NewWriter(os.Stdout)
			table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
			table.SetHeader([]string{tableHeaderModuleName, tableHeaderStatus, tableHeaderAddresses})
			for _, module := range arrayOfModules {
				addresses := ""
				connection := tableNotConnectedStatus
				instanceList := make([]string, 0, len(module.Status))
				for _, value := range module.Status {
					instanceList = append(instanceList, fmt.Sprintf("%s:%s %s", value.Address.IP, value.Address.Port, value.Version))
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
	return nil
}
