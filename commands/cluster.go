// Copyright 2015-2018 The elastic.go authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Author: Robin Hahling <robin.hahling@gw-computing.net>
// Author: Sergio Rubio <sergio@rubio.im>

// A command line tool to query the Elasticsearch REST API.

package commands

import (
	"encoding/json"
	"fmt"

	"strconv"

	"github.com/rubiojr/esg/types"
	"github.com/rubiojr/esg/util"
	"github.com/urfave/cli"
)

func init() {
	RegisterCommand(&cli.Command{
		Name:      "cluster",
		ShortName: "c",
		Usage:     "Get cluster information ",
		Subcommands: []cli.Command{
			{
				Name:      "health",
				ShortName: "he",
				Usage:     "Get cluster health",
				Action: func(c *cli.Context) {
					tunnel(c)
					out, err := getJSONBytes(cmdCluster(c, "health"))
					if err != nil {
						fatal(err)
					}
					printClusterHealth(out)
				},
			},
			{
				Name:      "state",
				ShortName: "s",
				Usage:     "Get cluster state",
				Action: func(c *cli.Context) {
					tunnel(c)
					out, err := getJSON(cmdCluster(c, "state"))
					if err != nil {
						fatal(err)
					}
					fmt.Println(out)
				},
			},
		},
	})
}

// command-line commands from now on
func cmdCluster(c *cli.Context, subCmd string) string {
	route := "_cluster/"
	url := c.GlobalString("baseurl")

	var arg string
	switch subCmd {
	case "health":
		arg = "health"
	case "state":
		arg = "state"
	default:
		arg = ""
	}
	return url + route + arg
}

func printClusterHealth(data []byte) {
	var b *types.ClusterHealth
	if err := json.Unmarshal(data, &b); err != nil {
		panic(err)
	}
	t := [][]string{
		[]string{"Name", b.ClusterName},
		[]string{"Status", b.Status},
		[]string{"Data Nodes", strconv.Itoa(b.NumberOfDataNodes)},
		[]string{"Pending Tasks", strconv.Itoa(b.NumberOfPendingTasks)},
		[]string{"In Flight Fetch", strconv.Itoa(b.NumberOfInFlightFetch)},
		[]string{"Task Max Waiting In Queue (ms)", strconv.Itoa(b.TaskMaxWaitingInQueueMillis)},
		[]string{"Active Primary Shards", strconv.Itoa(b.ActivePrimaryShards)},
		[]string{"Active Shards", strconv.Itoa(b.ActiveShards)},
		[]string{"Active Shards (%)", strconv.FormatFloat(b.ActiveShardsPercentAsNumber, 'f', -1, 32)},
		[]string{"Delayed Unassigned Shards", strconv.Itoa(b.DelayedUnassignedShards)},
		[]string{"Initializing Shards", strconv.Itoa(b.InitializingShards)},
		[]string{"Relocating Shards", strconv.Itoa(b.RelocatingShards)},
	}

	for _, v := range t {
		fmt.Printf("%s", util.PadRight(v[0], 40, " "))
		fmt.Printf("%s\n", v[1])
	}
}
