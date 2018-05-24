# A command-line tool to query the Elasticsearch REST API

Elasticsearch commandline interface.

Forked from https://github.com/Rolinh/elastic-go originally, aims to be a more modular with friendlier output than the original tool.

## Installation

Providing that [Go](https://golang.org) is installed and that `$GOPATH` is set,
simply use the following command:
```
go get -u github.com/rubiojr/esg
```

Make sure that `$GOPATH/bin` is in your `$PATH`.

## Usage

`esg help` provides general help:
```
$ esg help
NAME:
   esg - A command line tool to query the Elasticsearch REST API

USAGE:
   esg [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR(S):
   Robin Hahling <robin.hahling@gw-computing.net>

COMMANDS:
   cluster, c   Get cluster information
   index, i     Get index information
   node, n      Get cluster nodes information
   query, q     Perform any ES API GET query
   stats, s     Get statistics
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --baseurl "http://localhost:9200/"   Base API URL
   --help, -h                           show help
   --version, -v                        print the version
```

Help works for any subcommand as well. For instance:
```
$ esg index help
NAME:
   esg index - Get index information

USAGE:
   esg index [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   docs-count, dc       Get index documents count
   list, l              List all indexes
   size, si             Get index size
   status, st           Get index status
   verbose, v           List indexes information with many stats
   help, h              Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h   show help

```
