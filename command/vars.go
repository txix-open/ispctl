package command

import (
	"encoding/csv"
	"ispctl/bash"
	"ispctl/command/utils"
	"ispctl/model"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type VariablesService interface {
	GetAllVariables() ([]model.Variable, error)
	GetVariableByName(variableName string) (*model.Variable, error)
	DeleteVariable(variableName string) error
	UpsertVariables(vars []model.UpsertVariableRequest) error
}

type Variables struct {
	service      VariablesService
	autoComplete AutoComplete
}

func NewVariables(service VariablesService, autoComplete AutoComplete) Variables {
	return Variables{
		service:      service,
		autoComplete: autoComplete,
	}
}

func (c Variables) Command() *cli.Command {
	return &cli.Command{
		Name:  "vars",
		Usage: "manage variables",
		Subcommands: []*cli.Command{
			{
				Name:         "get",
				Usage:        "vars get <name>",
				Action:       c.getByName,
				BashComplete: c.autoComplete.Complete(bash.VariableName, bash.Empty),
			},
			{
				Name:   "list",
				Usage:  "list all variables",
				Action: c.list,
			},
			{
				Name:   "set",
				Usage:  "vars set <name> <value>",
				Action: c.set,
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: SecretFlagName, Usage: secretVariableFlagUsage},
				},
				BashComplete: c.autoComplete.Complete(bash.VariableName, bash.Empty),
			},
			{
				Name:         "delete",
				Usage:        "vars delete <name>",
				Action:       c.delete,
				BashComplete: c.autoComplete.Complete(bash.VariableName, bash.Empty),
			},
			{
				Name:   "upload",
				Usage:  "vars upload <filepath>",
				Action: c.upload,
			},
		},
	}
}

func (c Variables) list(ctx *cli.Context) error {
	vars, err := c.service.GetAllVariables()
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

func (c Variables) getByName(ctx *cli.Context) error {
	variableName := ctx.Args().Get(0)
	variable, err := c.service.GetVariableByName(variableName)
	if err != nil {
		return err
	}
	return utils.PrintAnswer(variable)
}

func (c Variables) set(ctx *cli.Context) error {
	variableName := ctx.Args().Get(0)
	variableValue := ctx.Args().Get(1)
	variableType := model.TextVariableType
	if ctx.Bool(SecretFlagName) {
		variableType = model.SecretVariableType
	}
	return c.service.UpsertVariables([]model.UpsertVariableRequest{{
		Name:  variableName,
		Value: variableValue,
		Type:  variableType,
	}})
}

func (c Variables) delete(ctx *cli.Context) error {
	return c.service.DeleteVariable(ctx.Args().Get(0))
}

func (c Variables) upload(ctx *cli.Context) error {
	variables, err := c.readVariablesFromCsv(ctx.Args().First())
	if err != nil {
		return err
	}
	return c.service.UpsertVariables(variables)
}

func (c Variables) readVariablesFromCsv(filepath string) ([]model.UpsertVariableRequest, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to open file %s", filepath)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read CSV file")
	}
	if len(records) < 2 {
		return nil, errors.New("CSV file must have a header and at least one data row")
	}

	header := records[0]
	columnIndexes := map[string]int{}
	for i, col := range header {
		columnIndexes[strings.ToLower(col)] = i
	}

	for _, columnName := range []string{"name", "type", "value", "description"} {
		_, ok := columnIndexes[columnName]
		if !ok {
			return nil, errors.Errorf("column '%s' is required", columnName)
		}
	}

	variables := make([]model.UpsertVariableRequest, 0, len(records)-1)
	for _, row := range records[1:] {
		variables = append(variables, model.UpsertVariableRequest{
			Name:        row[columnIndexes["name"]],
			Type:        row[columnIndexes["type"]],
			Value:       row[columnIndexes["value"]],
			Description: row[columnIndexes["description"]],
		})
	}
	return variables, nil
}
