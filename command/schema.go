package command

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	convert "github.com/integration-system/isp-lib/config/schema"
	"html/template"
	"isp-ctl/flag"
	"isp-ctl/service"
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
		BashComplete: bashSchema.run,
	}
}

var schema schemaCommand

type schemaCommand struct{}

func (s schemaCommand) action(c *cli.Context) {
	if err := checkFlags(c); err != nil {
		printError(err)
		return
	}
	moduleName := c.Args().First()

	if schemaConfig := s.getSchemaConfig(moduleName); schemaConfig != nil {
		schema := make(map[string]interface{})
		schema["title"] = moduleName
		schema["schema"] = schemaConfig
		switch c.String(flag.OutPrint.Name) {
		case flag.OutPrintJsonValue:
			printAnswer(schema)
		case flag.OutPrintHtmlValue:
			s.printHtml(schema)
		default:
			printError(errors.New(fmt.Sprintf(
				"invalid flag value, expected %s or %s", flag.OutPrintJsonValue, flag.OutPrintHtmlValue)))
		}
	}
}

func (s schemaCommand) getSchemaConfig(moduleName string) interface{} {
	if configuration, _ := getModuleConfiguration(moduleName); configuration == nil {
		return nil
	} else if schema, err := service.ConfigClient.GetSchemaByModuleId(configuration.ModuleId); err != nil {
		printError(err)
		return nil
	} else {
		return convert.DereferenceSchema(schema.Schema)
	}
}

func (s schemaCommand) printHtml(schema map[string]interface{}) {
	if temp, err := template.New("template").Parse(htmlTemplate); err != nil {
		printError(err)
		return
	} else {
		if err := temp.ExecuteTemplate(os.Stdout, "template", schema); err != nil {
			printError(err)
			return
		}
	}
}
