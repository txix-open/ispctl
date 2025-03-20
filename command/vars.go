package command

import (
	"encoding/csv"
	"fmt"
	"ispctl/bash"
	"ispctl/cfg"
	"ispctl/command/flag"
	"ispctl/command/utils"
	"ispctl/service"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

func VariablesCommands() *cli.Command {
	return &cli.Command{
		Name:  "vars",
		Usage: "manage variables",
		Subcommands: []*cli.Command{
			{
				Name:         "get",
				Usage:        "vars get <name>",
				Action:       varsComm.getByName,
				BashComplete: bash.Get(bash.VariableName, bash.Empty).Complete,
			},
			{
				Name:   "list",
				Usage:  "list all variables",
				Action: varsComm.list,
			},
			{
				Name:   "set",
				Usage:  "vars set <name> <value>",
				Action: varsComm.set,
				Flags: []cli.Flag{
					flag.SetVariableSecretType,
				},
				BashComplete: bash.Get(bash.VariableName, bash.Empty).Complete,
			},
			{
				Name:         "delete",
				Usage:        "vars delete <name>",
				Action:       varsComm.delete,
				BashComplete: bash.Get(bash.VariableName, bash.Empty).Complete,
			},
			{
				Name:   "upload",
				Usage:  "vars upload <filepath>",
				Action: varsComm.upload,
			},
		},
	}
}

var varsComm varsCommands

type varsCommands struct{}

func (g varsCommands) list(ctx *cli.Context) error {
	if err := flag.CheckGlobal(ctx); err != nil {
		return err
	}

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
	if err := flag.CheckGlobal(ctx); err != nil {
		return err
	}

	variableName := ctx.Args().Get(0)
	variable, err := service.Config.GetVariableByName(variableName)
	if err != nil {
		return err
	}
	utils.PrintAnswer(variable)
	return nil
}

func (g varsCommands) set(ctx *cli.Context) error {
	if err := flag.CheckGlobal(ctx); err != nil {
		return err
	}

	variableName := ctx.Args().Get(0)
	variableValue := ctx.Args().Get(1)
	variableType := cfg.TextVariableType

	if ctx.Bool(flag.SetVariableSecretType.Name) {
		variableType = cfg.SecretVariableType
	}
	return service.Config.UpsertVariables([]cfg.UpsertVariableRequest{{
		Name:  variableName,
		Value: variableValue,
		Type:  variableType,
	}})
}

func (g varsCommands) delete(ctx *cli.Context) error {
	if err := flag.CheckGlobal(ctx); err != nil {
		return err
	}

	variableName := ctx.Args().Get(0)
	return service.Config.DeleteVariable(variableName)
}

func (g varsCommands) upload(ctx *cli.Context) error {
	if err := flag.CheckGlobal(ctx); err != nil {
		return err
	}

	filepath := ctx.Args().First()
	variables, err := g.readVariablesFromCsv(filepath)
	if err != nil {
		return err
	}

	return service.Config.UpsertVariables(variables)
}

func (g varsCommands) readVariablesFromCsv(filepath string) ([]cfg.UpsertVariableRequest, error) {
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

	variables := make([]cfg.UpsertVariableRequest, 0, len(records)-1)
	for _, row := range records[1:] {
		variables = append(variables, cfg.UpsertVariableRequest{
			Name:        row[columnIndexes["name"]],
			Type:        row[columnIndexes["type"]],
			Value:       row[columnIndexes["value"]],
			Description: row[columnIndexes["description"]],
		})
	}
	return variables, nil
}
