package command

import (
	"fmt"
	"os"
	"strings"

	"ispctl/command/utils"
	"ispctl/model"

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

type StatusService interface {
	GetAvailableConfigs() ([]model.ModuleInfo, error)
}

type Status struct {
	service StatusService
}

func NewStatus(service StatusService) Status {
	return Status{service: service}
}

func (c Status) Command() *cli.Command {
	return &cli.Command{
		Name:   "status",
		Usage:  "get available configs",
		Action: c.action,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: OutPrintFlagName, Usage: outPrintStatusUsage},
		},
	}
}

func (c Status) action(ctx *cli.Context) error {
	arrayOfModules, err := c.service.GetAvailableConfigs()
	if err != nil {
		return err
	}
	switch ctx.String(OutPrintFlagName) {
	case OutPrintJsonValue:
		return utils.PrintAnswer(arrayOfModules)
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
	return nil
}
