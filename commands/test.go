package commands

import (
	"fmt"

	"github.com/urfave/cli"
)

func init() {
	cmd := cli.Command{
		Name:      "test",
		ShortName: "test",
		Usage:     "Test command",
		Subcommands: []cli.Command{
			{
				Name:      "hello",
				ShortName: "he",
				Usage:     "Test subcommand",
				Action: func(c *cli.Context) {
					fmt.Println("Hello world!")
				},
			},
		},
	}
	RegisterCommand(&cmd)
}
