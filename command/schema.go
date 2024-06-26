package command

import (
	"html/template"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"ispctl/bash"
	"ispctl/command/utils"
	"ispctl/flag"
	"ispctl/service"
	"ispctl/tmpl"
)

func Schema() *cli.Command {
	return &cli.Command{
		Name:   "schema",
		Usage:  "get schema configuration by module_name",
		Action: schema.action,
		Flags: []cli.Flag{
			flag.OutPrintSchema,
		},
		BashComplete: bash.Get(bash.ModuleName, bash.Empty).Complete,
	}
}

var schema schemaCommand

type schemaCommand struct{}

func (s schemaCommand) action(ctx *cli.Context) error {
	if err := flag.CheckGlobal(ctx); err != nil {
		return err
	}
	moduleName := ctx.Args().First()

	if schemaConfig := s.getSchemaConfig(moduleName); schemaConfig != nil {
		schema := make(map[string]interface{})
		schema["title"] = moduleName
		schema["schema"] = schemaConfig
		switch ctx.String(flag.OutPrintSchema.Name) {
		case flag.OutPrintJsonValue:
			utils.PrintAnswer(schema)
		case flag.OutPrintHtmlValue:
			s.printHtml(schema)
		default:
			return errors.Errorf("invalid flag value, expected [%s] or [%s]", flag.OutPrintJsonValue, flag.OutPrintHtmlValue)
		}
	}
	return nil
}

func (s schemaCommand) getSchemaConfig(moduleName string) interface{} {
	if configuration, err := service.Config.GetConfigurationByModuleName(moduleName); err != nil {
		utils.PrintError(err)
		return nil
	} else if schema, err := service.Config.GetSchemaByModuleId(configuration.ModuleId); err != nil {
		utils.PrintError(err)
		return nil
	} else {
		return schema
	}
}

func (s schemaCommand) printHtml(schema map[string]interface{}) {
	if temp, err := template.New("template").Parse(tmpl.HtmlFile); err != nil {
		utils.PrintError(err)
		return
	} else {
		if err := temp.ExecuteTemplate(os.Stdout, "template", schema); err != nil {
			utils.PrintError(err)
			return
		}
	}
}
