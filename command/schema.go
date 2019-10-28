package command

import (
	"github.com/codegangsta/cli"
	"github.com/pkg/errors"
	"html/template"
	"isp-ctl/bash"
	"isp-ctl/flag"
	"isp-ctl/service"
	"isp-ctl/tmpl"
	"os"
)

func Schema() cli.Command {
	return cli.Command{
		Name:   "schema",
		Usage:  "get schema configuration by module_name",
		Action: schema.action,
		Flags: []cli.Flag{
			flag.OutPrint,
		},
		BashComplete: bash.ModuleName.Complete,
	}
}

var schema schemaCommand

type schemaCommand struct{}

func (s schemaCommand) action(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		printError(err)
		return
	}
	moduleName := ctx.Args().First()

	if schemaConfig := s.getSchemaConfig(moduleName); schemaConfig != nil {
		schema := make(map[string]interface{})
		schema["title"] = moduleName
		schema["schema"] = schemaConfig
		switch ctx.String(flag.OutPrint.Name) {
		case flag.OutPrintJsonValue:
			printAnswer(schema)
		case flag.OutPrintHtmlValue:
			s.printHtml(schema)
		default:
			printError(errors.Errorf(
				"invalid flag value, expected %s or %s", flag.OutPrintJsonValue, flag.OutPrintHtmlValue))
		}
	}
}

func (s schemaCommand) getSchemaConfig(moduleName string) interface{} {
	if configuration, _, err := service.Config.GetConfigurationAndJsonByModuleName(moduleName); err != nil {
		printError(err)
		return nil
	} else if schema, err := service.Config.GetSchemaByModuleId(configuration.ModuleId); err != nil {
		printError(err)
		return nil
	} else {
		return schema
	}
}

func (s schemaCommand) printHtml(schema map[string]interface{}) {
	if temp, err := template.New("template").Parse(tmpl.HtmlFile); err != nil {
		printError(err)
		return
	} else {
		if err := temp.ExecuteTemplate(os.Stdout, "template", schema); err != nil {
			printError(err)
			return
		}
	}
}
