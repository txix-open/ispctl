package command

import (
	"html/template"
	"os"

	"ispctl/bash"
	"ispctl/command/utils"
	"ispctl/model"
	"ispctl/tmpl"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/rc/schema"
	"github.com/urfave/cli/v2"
)

type SchemaService interface {
	GetConfigurationByModuleName(moduleName string) (*model.Config, error)
	GetSchemaByModuleId(moduleId string) (schema.Schema, error)
}

type Schema struct {
	service      SchemaService
	autoComplete AutoComplete
}

func NewSchema(service SchemaService, autoComplete AutoComplete) Schema {
	return Schema{
		service:      service,
		autoComplete: autoComplete,
	}
}

func (c Schema) Command() *cli.Command {
	return &cli.Command{
		Name:   "schema",
		Usage:  "get schema configuration by module_name",
		Action: c.action,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: OutPrintFlagName, Usage: outPrintSchemaUsage, Value: OutPrintJsonValue},
		},
		BashComplete: c.autoComplete.Complete(bash.ModuleName, bash.Empty),
	}
}

func (c Schema) action(ctx *cli.Context) error {
	moduleName := ctx.Args().First()
	schemaConfig, err := c.getSchemaConfig(moduleName)
	if err != nil {
		return err
	}

	schema := make(map[string]any)
	schema["title"] = moduleName
	schema["schema"] = schemaConfig
	switch ctx.String(OutPrintFlagName) {
	case OutPrintJsonValue:
		return utils.PrintAnswer(schema)
	case OutPrintHtmlValue:
		return c.printHtml(schema)
	default:
		return errors.Errorf("invalid flag value, expected [%s] or [%s]", OutPrintJsonValue, OutPrintHtmlValue)
	}
}

func (c Schema) getSchemaConfig(moduleName string) (schema.Schema, error) {
	configuration, err := c.service.GetConfigurationByModuleName(moduleName)
	if err != nil {
		return nil, err
	}
	schema, err := c.service.GetSchemaByModuleId(configuration.ModuleId)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

func (c Schema) printHtml(schema map[string]any) error {
	temp, err := template.New("template").Parse(tmpl.HtmlFile)
	if err != nil {
		return err
	}
	return temp.ExecuteTemplate(os.Stdout, "template", schema)
}
