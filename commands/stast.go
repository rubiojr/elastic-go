// Copyright 2015-2018 The elastic.go authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Author: Robin Hahling <robin.hahling@gw-computing.net>
// Author: Sergio Rubio <sergio@rubio.im>
package commands

import (
	"fmt"

	"github.com/urfave/cli"
)

func init() {
	RegisterCommand(
		&cli.Command{
			Name:      "stats",
			ShortName: "s",
			Usage:     "Get statistics",
			Subcommands: []cli.Command{
				{
					Name:      "size",
					ShortName: "s",
					Usage:     "Get index sizes",
					Action: func(c *cli.Context) {
						tunnel(c)
						out, err := getJSON(cmdStats(c, "size"))
						if err != nil {
							fatal(err)
						}
						fmt.Println(out)
					},
				},
			},
		},
	)
}

func cmdStats(c *cli.Context, subCmd string) string {
	var route string
	url := c.GlobalString("baseurl")
	switch subCmd {
	case "size":
		route = "_stats/index,store"
	default:
		route = ""
	}
	return url + route
}
