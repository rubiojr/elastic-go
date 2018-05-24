// Copyright 2015-2018 The elastic.go authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Author: Robin Hahling <robin.hahling@gw-computing.net>
// Author: Sergio Rubio <sergio@rubio.im>

package commands

import (
	"bufio"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func init() {
	RegisterCommand(&cli.Command{
		Name:      "index",
		ShortName: "i",
		Usage:     "Get index information",
		Subcommands: []cli.Command{
			{
				Name:      "docs-count",
				ShortName: "dc",
				Usage:     "Get index documents count",
				Action: func(c *cli.Context) {
					tunnel(c)
					list, err := getRaw(cmdIndex(c, "list"))
					if err != nil {
						fatal(err)
					}
					for _, idx := range filteredDocsCountIndexes(list) {
						fmt.Println(idx)
					}
				},
			},
			{
				Name:      "list",
				ShortName: "l",
				Usage:     "List all indexes",
				Action: func(c *cli.Context) {
					tunnel(c)
					list, err := getRaw(cmdIndex(c, "list"))
					if err != nil {
						fatal(err)
					}
					for _, idx := range filteredListIndexes(list) {
						fmt.Println(idx)
					}
				},
			},
			{
				Name:      "size",
				ShortName: "si",
				Usage:     "Get index size",
				Action: func(c *cli.Context) {
					tunnel(c)
					list, err := getRaw(cmdIndex(c, "list"))
					if err != nil {
						fatal(err)
					}
					for _, idx := range filteredSizeIndexes(list) {
						fmt.Println(idx)
					}
				},
			},
			{
				Name:      "status",
				ShortName: "st",
				Usage:     "Get index status",
				Action: func(c *cli.Context) {
					tunnel(c)
					list, err := getRaw(cmdIndex(c, "list"))
					if err != nil {
						fatal(err)
					}
					for _, idx := range filteredStatusIndexes(list) {
						fmt.Println(idx)
					}
				},
			},
			{
				Name:      "verbose",
				ShortName: "v",
				Usage:     "List indexes information with many stats",
				Action: func(c *cli.Context) {
					tunnel(c)
					list, err := getRaw(cmdIndex(c, "list"))
					if err != nil {
						fatal(err)
					}
					fmt.Println(list)
				},
			},
		},
	})
}

// processing functions
func filteredDocsCountIndexes(list string) []string {
	var out []string
	scanner := bufio.NewScanner(strings.NewReader(list))
	counter := 0
	for scanner.Scan() {
		elmts := strings.Fields(scanner.Text())
		if len(elmts) < 6 || counter == 0 {
			counter += 1
			continue
		}
		out = append(out, fmt.Sprintf("%10s %s", colorizeStatus(elmts[5]), elmts[2]))
	}
	sort.Strings(out)
	return out
}

func filteredListIndexes(list string) []string {
	var out []string
	scanner := bufio.NewScanner(strings.NewReader(list))
	for scanner.Scan() {
		elmts := strings.Fields(scanner.Text())
		if len(elmts) < 3 {
			continue
		}
		out = append(out, elmts[2])
	}
	sort.Strings(out)

	return out
}

func filteredStatusIndexes(list string) []string {
	var out []string
	scanner := bufio.NewScanner(strings.NewReader(list))
	for scanner.Scan() {
		elmts := strings.Fields(scanner.Text())
		if len(elmts) < 3 {
			continue
		}
		out = append(out, fmt.Sprintf("%22s %s", colorizeStatus(elmts[0]), elmts[2]))
	}
	return out
}

func filteredSizeIndexes(list string) []string {
	var out []string
	scanner := bufio.NewScanner(strings.NewReader(list))
	for scanner.Scan() {
		elmts := strings.Fields(scanner.Text())
		if len(elmts) < 8 {
			continue
		}
		// ES 6.X output lists UUIDs for indices also, so it has one more element
		if len(elmts) < 10 {
			out = append(out, fmt.Sprintf("%10s %s", elmts[7], elmts[2]))
		} else {
			out = append(out, fmt.Sprintf("%10s %s", elmts[8], elmts[2]))
		}
	}

	start := time.Now()

	sort.Strings(out)

	mb := regexp.MustCompile(`[0-9\\.]+mb`)
	gb := regexp.MustCompile(`[0-9\\.]+gb`)
	kb := regexp.MustCompile(`[0-9\\.]+kb`)
	m := make(map[string][]string)
	for _, v := range out {
		if mb.MatchString(v) {
			m["mb"] = append(m["mb"], v)
		} else if gb.MatchString(v) {
			m["gb"] = append(m["gb"], v)
		} else if kb.MatchString(v) {
			m["kb"] = append(m["kb"], v)
		}
	}
	na := []string{}
	na = append(na, m["kb"]...)
	na = append(na, m["mb"]...)
	na = append(na, m["gb"]...)

	elapsed := time.Since(start)
	log.Debugf("Sorting took %s", elapsed)

	return na
}

func cmdIndex(c *cli.Context, subCmd string) string {
	var route string
	url := c.GlobalString("baseurl")
	switch subCmd {
	case "list":
		route = "_cat/indices?v"
	default:
		route = ""
	}
	return url + route
}
