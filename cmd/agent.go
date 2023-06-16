package cmd

import (
	"fmt"
	"github.com/VintageOps/dns-zone-compare/pkg/zonecompare"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

type opts struct {
	domain           string
	origin           string
	ignoreTTL        bool
	ignore           []string
	deep             []string
	deepAll          bool
	found            bool
	notfound         bool
	strict           bool
	json             bool
	text             bool
	countText        int
	destination      string
	labelorigin      string
	labeldestination string
}

func fatalOnErr(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func Execute() {
	// TODO: add json or text only outputs
	options := zonecompare.Opts{}
	app := &cli.App{
		Name:      "zonecompare",
		Usage:     "compare two dns zone files",
		ArgsUsage: "<path_zonefile1>|<address|name:port> <path_zonefile2>|<address|name:port>",
		Description: "zonecompare reads or transfer two DNS zone files and by default, output the differences.\n" +
			"It can also output the similarities and has an extensive set of options to customize the comparison.\n" +
			"The default output is a timestamped text highlighting the differences (or similarities with <--showfound>/<-f>),\n" +
			"But with the appropriate <--json>/<-j> option, it can print the same in a comprehensive json format.\n" +
			"The mandatory zonefiles argument can be specified either using their file path (e.g. tmp/zonefile1 tmp/zonefile2\n" +
			"Or using <address|name>:<port> format (e.g. localhost:8053 192.168.0.1:53), in which case, for a DNS zone transfer(axfr)",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "ignorettl",
				Value:       false,
				Aliases:     []string{"t"},
				Usage:       "Force TTL value to 604800 in both zones",
				Destination: &options.IgnoreTTL,
			},
			&cli.StringFlag{
				Name: "domain",
				Usage: "domain to compare (e.g. example.com)." +
					"Required when arguments are <ip|name>:<port>, " +
					"and on zonefiles with no $ORIGIN",
				Destination: &options.Domain,
			},
			&cli.StringFlag{
				Name:        "labelorigin",
				Aliases:     []string{"lo"},
				Usage:       "label of the origin zone, default origin filename|server:port",
				Destination: &options.Labelorigin,
			},
			&cli.StringFlag{
				Name:        "labeldestination",
				Aliases:     []string{"ld"},
				Usage:       "label of the destination zone, default destination filename|server:port",
				Destination: &options.Labeldestination,
			},
			&cli.BoolFlag{
				Name:        "showfound",
				Aliases:     []string{"f"},
				Value:       false,
				Usage:       "Report on found records",
				Destination: &options.Found,
			},
			&cli.BoolFlag{
				Name:        "skipnotfound",
				Aliases:     []string{"n"},
				Value:       false,
				Usage:       "Skip not found records",
				Destination: &options.Notfound,
			},
			&cli.BoolFlag{
				Name:        "strict",
				Value:       false,
				Aliases:     []string{"s"},
				Usage:       "Consider the different order of the same record a difference",
				Destination: &options.Strict,
			},
			&cli.StringSliceFlag{
				Name:    "ignore",
				Aliases: []string{"i"},
				Usage:   "Ignore <value> type records",
			},
			&cli.StringSliceFlag{
				Name:    "deep",
				Aliases: []string{"d"},
				Usage:   "Inspect <value> type records by merging, then splitting and sorting the content",
			},
			&cli.BoolFlag{
				Name:        "deepAll",
				Value:       false,
				Aliases:     []string{"da"},
				Usage:       "Inspect all type records by merging, then splitting and sorting the content",
				Destination: &options.DeepAll,
			},
			&cli.BoolFlag{
				Name:        "json",
				Value:       false,
				Aliases:     []string{"j"},
				Usage:       "output in json format",
				Destination: &options.Json,
			},
			&cli.BoolFlag{
				Name:        "text",
				Aliases:     []string{"x"},
				Usage:       "Forcing timestamped text in output, useful only to produce both json and text output",
				Destination: &options.Text,
				Value:       false,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				fmt.Printf("ERROR: Needs two arguments, the arguments being either the zonefiles path, or ip/name:<port>, %d provided.\n\n", c.NArg())
				cli.ShowAppHelpAndExit(c, 1)
			}
			options.Origin = c.Args().Get(0)
			options.Destination = c.Args().Get(1)
			options.CountText = c.Count("text")
			if options.Labelorigin == "" {
				options.Labelorigin = options.Origin
			}
			if options.Labeldestination == "" {
				options.Labeldestination = options.Destination
			}
			options.Ignore = c.StringSlice("ignore")
			options.Deep = c.StringSlice("deep")
			json_output := zonecompare.ZoneCompare(options)
			if options.Json {
				fmt.Println(json_output)
			}
			return cli.Exit("", 0)
		},
	}
	err := app.Run(os.Args)
	fatalOnErr(err)
}
