package commands

import (
	"encoding/json"
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/rubiojr/elastic-go/util"
	"github.com/urfave/cli"
)

func init() {
	RegisterCommand(&cli.Command{
		Name:      "node",
		ShortName: "n",
		Usage:     "Get cluster nodes information",
		Subcommands: []cli.Command{
			{
				Name:      "list",
				ShortName: "l",
				Usage:     "List nodes information",
				Action: func(c *cli.Context) {
					tunnel(c)
					out, err := getJSONBytes(cmdNode(c, "list"))
					if err != nil {
						fatal(err)
					}
					listNodes(out)
				},
			},
		},
	},
	)
}

func cmdNode(c *cli.Context, subCmd string) string {
	var route string
	url := c.GlobalString("baseurl")
	switch subCmd {
	case "list":
		route = "_nodes/_all/host,ip"
	default:
		route = ""
	}
	return url + route
}

func listNodes(data []byte) {
	var j map[string]interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		panic(err)
	}

	jsonParsed, err := gabs.ParseJSON(data)
	if err != nil {
		panic(err)
	}

	children, err := jsonParsed.S("nodes").ChildrenMap()
	if err != nil {
		panic(err)
	}

	for _, v := range children {
		t := [][]string{
			[]string{"Name", v.Path("name").String()},
			[]string{"Host", v.Path("host").String()},
			[]string{"IP", v.Path("ip").String()},
			[]string{"Version", v.Path("version").String()},
			[]string{"Build", v.Path("build").String()},
			[]string{"Master", v.Path("attributes.master").String()},
			[]string{"Role", v.Path("attributes.role").String()},
			[]string{"", ""},
		}
		for _, n := range t {
			fmt.Printf("%s", util.PadRight(n[0], 40, " "))
			fmt.Printf("%s\n", n[1])
		}
	}

}
