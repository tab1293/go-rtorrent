package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mrobinsn/go-rtorrent/rtorrent"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var (
	name    = "rTorrent XMLRPC CLI"
	version = "1.0.0"
	app     = initApp()
	conn    *rtorrent.RTorrent

	endpoint         string
	view             string
	hash             string
	torrentPath      string
	torrentURL       string
	fileIndex        int
	filePriority     int
	disableCertCheck bool
)

func initApp() *cli.App {
	nApp := cli.NewApp()

	nApp.Name = name
	nApp.Version = version
	nApp.Authors = []cli.Author{
		{Name: "Michael Robinson", Email: "m@michaelrobinson.io"},
	}

	// Global flags
	nApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "endpoint",
			Usage:       "rTorrent endpoint",
			Value:       "http://myrtorrent/RPC2",
			Destination: &endpoint,
		},
		cli.BoolFlag{
			Name:        "disable-cert-check",
			Usage:       "disable certificate checking on this endpoint, useful for testing",
			Destination: &disableCertCheck,
		},
	}

	nApp.Before = setupConnection

	nApp.Commands = []cli.Command{{
		Name:   "get-ip",
		Usage:  "retrieves the IP for this rTorrent instance",
		Action: getIP,
	}, {
		Name:   "get-name",
		Usage:  "retrieves the name for this rTorrent instance",
		Action: getName,
	}, {
		Name:   "get-totals",
		Usage:  "retrieves the up/down totals for this rTorrent instance",
		Action: getTotals,
	}, {
		Name:   "get-torrents",
		Usage:  "retrieves the torrents from this rTorrent instance",
		Action: getTorrents,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "view",
				Usage:       "view to use, known values: main, started, stopped, hashing, seeding",
				Value:       string(rtorrent.ViewMain),
				Destination: &view,
			},
		},
	}, {
		Name:   "get-files",
		Usage:  "retrieves the files for a specific torrent",
		Action: getFiles,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "hash",
				Usage:       "hash of the torrent",
				Value:       "unknown",
				Destination: &hash,
			},
		},
	}, {
		Name:   "start-torrent",
		Usage:  "starts a specific torrent",
		Action: startTorrent,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "hash",
				Usage:       "hash of the torrent",
				Value:       "unknown",
				Destination: &hash,
			},
		},
	}, {
		Name:   "stop-torrent",
		Usage:  "stops a specific torrent",
		Action: stopTorrent,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "hash",
				Usage:       "hash of the torrent",
				Value:       "unknown",
				Destination: &hash,
			},
		},
	}, {
		Name:   "add-torrent",
		Usage:  "add torrent from file",
		Action: addTorrentFile,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "file",
				Usage:       "path to torrent file",
				Value:       "unknown",
				Destination: &torrentPath,
			},
		},
	}, {
		Name:   "add-torrent-url",
		Usage:  "add torrent from URL",
		Action: addTorrentURL,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "url",
				Usage:       "url of the torrent",
				Value:       "unknown",
				Destination: &torrentURL,
			},
		},
	}, {
		Name:   "close-torrent",
		Usage:  "close torrent with hash",
		Action: closeTorrent,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "hash",
				Usage:       "hash of the torrent",
				Value:       "unknown",
				Destination: &hash,
			},
		},
	}, {
		Name:   "delete-torrent",
		Usage:  "delete torrent with hash",
		Action: deleteTorrent,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "hash",
				Usage:       "hash of the torrent",
				Value:       "unknown",
				Destination: &hash,
			},
		},
	}, {
		Name:   "get-torrent",
		Usage:  "get torrent with hash",
		Action: getTorrent,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "hash",
				Usage:       "hash of the torrent",
				Value:       "unknown",
				Destination: &hash,
			},
		},
	}, {
		Name:   "list-methods",
		Usage:  "list rpc methods",
		Action: listMethods,
	}, {
		Name:   "shutdown",
		Usage:  "shutdown",
		Action: shutdown,
	}, {
		Name:   "set-file-priority",
		Usage:  "Set a file priority",
		Action: setFilePriority,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "hash",
				Usage:       "hash of the torrent",
				Value:       "unknown",
				Destination: &hash,
			},
			cli.IntFlag{
				Name:        "index",
				Usage:       "index of the file in the torrent",
				Value:       -1,
				Destination: &fileIndex,
			},
			cli.IntFlag{
				Name:        "priority",
				Usage:       "index of the file in the torrent",
				Value:       0,
				Destination: &filePriority,
			},
		},
	}, {
		Name:   "skip-all-files",
		Usage:  "Set skip priority for all files",
		Action: skipAllFiles,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "hash",
				Usage:       "hash of the torrent",
				Value:       "unknown",
				Destination: &hash,
			},
		},
	},
	}

	return nApp
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func setupConnection(c *cli.Context) error {
	if endpoint == "" {
		return errors.New("endpoint must be specified")
	}
	conn = rtorrent.New(endpoint, disableCertCheck)
	return nil
}

func getIP(c *cli.Context) error {
	ip, err := conn.IP()
	if err != nil {
		return errors.Wrap(err, "failed to get rTorrent IP")
	}
	fmt.Println(ip)
	return nil
}

func getName(c *cli.Context) error {
	name, err := conn.Name()
	if err != nil {
		return errors.Wrap(err, "failed to get rTorrent name")
	}
	fmt.Println(name)
	return nil
}

func getTotals(c *cli.Context) error {
	// Get Down Total
	downTotal, err := conn.DownTotal()
	if err != nil {
		return errors.Wrap(err, "failed to get rTorrent down total")
	}
	fmt.Printf("%d\n", downTotal)

	// Get Up Total
	upTotal, err := conn.UpTotal()
	if err != nil {
		return errors.Wrap(err, "failed to get rTorrent up total")
	}
	fmt.Printf("%d\n", upTotal)
	return nil
}

func getTorrents(c *cli.Context) error {
	torrents, err := conn.GetTorrents(rtorrent.View(view))
	if err != nil {
		return errors.Wrap(err, "failed to get torrents")
	}
	for _, torrent := range torrents {
		fmt.Println(torrent.Pretty())
	}
	return nil
}

func getFiles(c *cli.Context) error {
	files, err := conn.GetFiles(rtorrent.Torrent{Hash: hash})
	if err != nil {
		return errors.Wrap(err, "failed to get files")
	}
	for _, file := range files {
		fmt.Println(file.Pretty())
	}
	return nil
}

func getTorrent(c *cli.Context) error {
	t, err := conn.GetTorrent(rtorrent.Torrent{Hash: hash})
	if err != nil {
		return errors.Wrap(err, "failed to get torrent")
	}

	fmt.Println(t.Pretty())

	return nil
}

func addTorrentFile(c *cli.Context) error {
	b, err := ioutil.ReadFile(torrentPath)
	if err != nil {
		return err
	}
	err = conn.AddTorrent(b)
	if err != nil {
		return errors.Wrap(err, "failed to start torrent")
	}

	return nil
}

func startTorrent(c *cli.Context) error {
	err := conn.StartTorrent(rtorrent.Torrent{Hash: hash})
	if err != nil {
		return errors.Wrap(err, "failed to start torrent")
	}

	return nil
}

func stopTorrent(c *cli.Context) error {
	err := conn.StartTorrent(rtorrent.Torrent{Hash: hash})
	if err != nil {
		return errors.Wrap(err, "failed to start torrent")
	}

	return nil
}

func closeTorrent(c *cli.Context) error {
	err := conn.CloseTorrent(rtorrent.Torrent{Hash: hash})
	if err != nil {
		return errors.Wrap(err, "failed to start torrent")
	}

	return nil
}

func deleteTorrent(c *cli.Context) error {
	err := conn.Delete(rtorrent.Torrent{Hash: hash})
	if err != nil {
		return errors.Wrap(err, "failed to start torrent")
	}

	return nil
}

func addTorrentURL(c *cli.Context) error {
	fmt.Printf("torrent url %s\n", torrentURL)
	err := conn.AddTorrentURL(torrentURL)
	if err != nil {
		return errors.Wrap(err, "failed to add torrent")
	}

	return nil
}

func listMethods(c *cli.Context) error {
	methods, err := conn.ListMethods()
	if err != nil {
		return errors.Wrap(err, "failed to list methods")
	}

	for _, method := range methods {
		fmt.Println(method)
	}

	return nil
}

func shutdown(c *cli.Context) error {
	err := conn.Shutdown()
	if err != nil {
		return errors.Wrap(err, "failed to list methods")
	}

	return nil
}

func setFilePriority(c *cli.Context) error {
	err := conn.SetFilePrority(rtorrent.Torrent{Hash: hash}, fileIndex, filePriority)
	if err != nil {
		return errors.Wrap(err, "failed to list methods")
	}

	return nil
}

func skipAllFiles(c *cli.Context) error {
	files, err := conn.GetFiles(rtorrent.Torrent{Hash: hash})
	if err != nil {
		return errors.Wrap(err, "failed to get files")
	}
	for i := range files {
		err := conn.SetFilePrority(rtorrent.Torrent{Hash: hash}, i, 0)
		if err != nil {
			return errors.Wrap(err, "failed to list methods")
		}
	}
	return nil

}
