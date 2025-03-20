package command

import "github.com/urfave/cli/v2"

func AllCommands() []*cli.Command {
	return []*cli.Command{
		Status(),
		Get(),
		Set(),
		Delete(),
		Schema(),
		Merge(),
		GitGet(),
		VariablesCommands(),
	}
}
