package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rubiojr/elastic-go/sshtunnel"
	"github.com/urfave/cli"

	"github.com/gilliek/go-xterm256/xterm256"
	prettyjson "github.com/hokaccha/go-prettyjson"
)

var _tunnel *exec.Cmd
var app *cli.App = cli.NewApp()

func RegisterCommand(command *cli.Command) {
	app.Commands = append(app.Commands, *command)
}

func Init() {
	app.Name = "elastic"
	app.Usage = "A command line tool to query the Elasticsearch REST API"
	app.Version = "1.0.1"
	app.Author = "Robin Hahling"
	app.Email = "robin.hahling@gw-computing.net"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "baseurl",
			Value: "http://localhost:9200/",
			Usage: "Base API URL",
		},
		cli.StringFlag{
			Name:  "tunnel-host",
			Usage: "Tunnel host",
		},
		cli.StringFlag{
			Name:  "tunnel-user",
			Usage: "Tunnel user",
		},
		cli.IntFlag{
			Name:  "tunnel-port",
			Usage: "Tunnel port",
			Value: 22,
		},
		cli.StringFlag{
			Name:  "tunnel-endpoint",
			Usage: "Tunnel endpoint",
			Value: "9200:localhost:9200",
		},
	}
	defer func() {
		if _tunnel != nil {
			_tunnel.Process.Kill()
		}
	}()
	app.Run(os.Args)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func getJSONBytes(route string) ([]byte, error) {
	r, err := http.Get(route)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %s", r.Status)
	}

	mediatype, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	if mediatype == "" {
		return nil, errors.New("mediatype not set")
	}
	if mediatype != "application/json" {
		return nil, fmt.Errorf("mediatype is '%s', 'application/json' expected", mediatype)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		return data, nil
	}
	return nil, err
}

func getJSON(route string) (string, error) {
	r, err := http.Get(route)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code: %s", r.Status)
	}

	mediatype, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}
	if mediatype == "" {
		return "", errors.New("mediatype not set")
	}
	if mediatype != "application/json" {
		return "", fmt.Errorf("mediatype is '%s', 'application/json' expected", mediatype)
	}

	var b interface{}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		return "", err
	}

	out, err := prettyjson.Marshal(b)
	return string(out), err
}

func getRaw(route string) (string, error) {
	r, err := http.Get(route)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code: %s", r.Status)
	}

	out, err := ioutil.ReadAll(r.Body)
	return string(out), err
}

func colorizeStatus(status string) string {
	var color xterm256.Color
	switch status {
	case "red":
		color = xterm256.Red
	case "green":
		color = xterm256.Green
	case "yellow":
		color = xterm256.Yellow
	default:
		return status
	}
	return xterm256.Sprint(color, status)
}

func tunnel(c *cli.Context) {
	if c.GlobalString("tunnel-host") != "" {
		endpoint := c.GlobalString("tunnel-endpoint")
		_tunnel = sshtunnel.Tunnel(c.GlobalString("tunnel-user"), c.GlobalString("tunnel-host"), c.GlobalInt("tunnel-port"), endpoint)
		hostport := strings.Split(endpoint, ":")[1:3]
		for i := 1; i <= 10; i++ {
			conn, err := net.Dial("tcp", strings.Join(hostport, ":"))
			if err == nil {
				defer conn.Close()
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
