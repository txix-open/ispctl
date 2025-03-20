package command

import (
	"encoding/csv"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	"ispctl/bash"
	"ispctl/command/flag"
	"ispctl/command/utils"
	"ispctl/entity"
	"ispctl/service"
	"os"
	"strings"
)

func VariablesCommands() *cli.Command {
	return &cli.Command{
		Name:  "vars",
		Usage: "manage variables",
		Subcommands: []*cli.Command{
			{
				Name:         "get",
				Usage:        "vars get <name>",
				Before:       flag.CheckGlobal,
				Action:       varsComm.getByName,
				BashComplete: bash.Get(bash.VariableName, bash.Empty).Complete,
			},
			{
				Name:   "list",
				Usage:  "list all variables",
				Before: flag.CheckGlobal,
				Action: varsComm.list,
			},
			{
				Name:   "set",
				Usage:  "vars set <name> <value>",
				Before: flag.CheckGlobal,
				Action: varsComm.set,
				Flags: []cli.Flag{
					flag.SetVariableSecretType,
				},
				BashComplete: bash.Get(bash.VariableName, bash.Empty).Complete,
			},
			{
				Name:         "delete",
				Usage:        "vars delete <name>",
				Before:       flag.CheckGlobal,
				Action:       varsComm.delete,
				BashComplete: bash.Get(bash.VariableName, bash.Empty).Complete,
			},
			{
				Name:   "upload",
				Usage:  "vars upload <filepath>",
				Before: flag.CheckGlobal,
				Action: varsComm.upload,
			},
		},
	}
}

var varsComm varsCommands

type varsCommands struct{}

func (g varsCommands) list(ctx *cli.Context) error {
	vars, err := service.Config.GetAllVariables()
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description", "Type", "Value", "Configs"})
	for _, v := range vars {
		configs := make([]string, 0, len(v.ContainsInConfigs))
		for _, c := range v.ContainsInConfigs {
			configs = append(configs, c.Name)
		}
		table.Append([]string{v.Name, v.Description, v.Type, v.Value, strings.Join(configs, ",")})
	}
	table.SetBorder(true)
	table.SetRowLine(true)
	table.Render()
	return nil
}

func (g varsCommands) getByName(ctx *cli.Context) error {
	variableName := ctx.Args().Get(0)
	variable, err := service.Config.GetVariableByName(variableName)
	if err != nil {
		return err
	}
	utils.PrintAnswer(variable)
	return nil
}

func (g varsCommands) set(ctx *cli.Context) error {
	variableName := ctx.Args().Get(0)
	variableValue := ctx.Args().Get(1)
	variableType := entity.TextVariableType
	if ctx.Bool(flag.SetVariableSecretType.Name) {
		variableType = entity.SecretVariableType
	}
	return service.Config.UpsertVariables([]entity.UpsertVariableRequest{{
		Name:  variableName,
		Value: variableValue,
		Type:  variableType,
	}})
}

func (g varsCommands) delete(ctx *cli.Context) error {
	return service.Config.DeleteVariable(ctx.Args().Get(0))
}

func (g varsCommands) upload(ctx *cli.Context) error {
	variables, err := g.readVariablesFromCsv(ctx.Args().First())
	if err != nil {
		return err
	}
	return service.Config.UpsertVariables(variables)
}

func (g varsCommands) readVariablesFromCsv(filepath string) ([]entity.UpsertVariableRequest, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file must have a header and at least one data row")
	}
	header := records[0]
	columnIndexes := map[string]int{}
	for i, col := range header {
		columnIndexes[strings.ToLower(col)] = i
	}
	variables := make([]entity.UpsertVariableRequest, 0, len(records)-1)
	for _, row := range records[1:] {
		variables = append(variables, entity.UpsertVariableRequest{
			Name:        row[columnIndexes["name"]],
			Type:        row[columnIndexes["type"]],
			Value:       row[columnIndexes["value"]],
			Description: row[columnIndexes["description"]],
		})
	}
	return variables, nil
}
