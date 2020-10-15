# go-rtorrent
[![GoDoc](https://godoc.org/github.com/tab1293/go-rtorrent/rtorrent?status.svg)](https://godoc.org/github.com/tab1293/go-rtorrent/rtorrent)
[![Go Report Card](https://goreportcard.com/badge/github.com/tab1293/go-rtorrent)](https://goreportcard.com/report/github.com/tab1293/go-rtorrent)
[![Build Status](https://travis-ci.org/tab1293/go-rtorrent.svg?branch=master)](https://travis-ci.org/tab1293/go-rtorrent)
[![Coverage Status](https://coveralls.io/repos/github/tab1293/go-rtorrent/badge.svg?branch=master)](https://coveralls.io/github/tab1293/go-rtorrent?branch=master)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)


> rTorrent XMLRPC Bindings for Go (golang)

## Documentation
[GoDoc](https://godoc.org/github.com/tab1293/go-rtorrent/rtorrent)

## Features
- Get IP, Name, Up/Down totals
- Get torrents within a view
- Get torrent by hash
- Get files for torrents
- Set the label on a torrent
- Add a torrent by URL or by metadata
- Delete a torrent (including files)

## Installation
To install the package, run `go get github.com/tab1293/go-rtorrent`

To use it in application, import `"github.com/tab1293/go-rtorrent/rtorrent"`

To install the command line utility, run `go install "github.com/tab1293/go-rtorrent"`

## Library Usage

```
conn, _ := rtorrent.New("http://my-rtorrent.com/RPC2", false)
name, _ := conn.Name()
fmt.Printf("My rTorrent's name: %v", name)
```

You can connect to a server using Basic Authentication by including the credentials in the endpoint URL:
```
conn, _ := rtorrent.New("https://user:pass@my-rtorrent.com/RPC2", false)
```

## Command Line Utility
A basic command line utility is included

`go-rtorrent`

```
NAME:
   rTorrent XMLRPC CLI - A new cli application

USAGE:
   go-rtorrent [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR(S):
   Michael Robinson <m@michaelrobinson.io>

COMMANDS:
   get-ip    retrieves the IP for this rTorrent instance
   get-name    retrieves the name for this rTorrent instance
   get-totals    retrieves the up/down totals for this rTorrent instance
   get-torrents    retrieves the torrents from this rTorrent instance
   get-files    retrieves the files for a specific torrent
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --endpoint "http://myrtorrent/RPC2"    rTorrent endpoint
   --disable-cert-check            disable certificate checking on this endpoint, useful for testing
   --help, -h                show help
   --version, -v            print the version
```

## Contributing
Pull requests are welcome, please ensure you add relevant tests for any new/changed functionality.
