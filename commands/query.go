// Copyright 2015-2018 The elastic.go authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Author: Robin Hahling <robin.hahling@gw-computing.net>
// Author: Sergio Rubio <sergio@rubio.im>
package commands

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

func init() {
	RegisterCommand(
		&cli.Command{
			Name:      "query",
			ShortName: "q",
			Usage:     "Perform any ES API GET query",
			Action: func(c *cli.Context) {
				tunnel(c)
				var out string
				var err error
				if strings.Contains(c.Args().First(), "_cat/") {
					out, err = getRaw(cmdQuery(c))
				} else {
					out, err = getJSON(cmdQuery(c))
				}
				if err != nil {
					fatal(err)
				}
				fmt.Println(out)
			},
		},
	)
}

func cmdQuery(c *cli.Context) string {
	route := c.Args().First()
	url := c.GlobalString("baseurl")
	return url + route
}
