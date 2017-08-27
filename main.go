package main

import (
	"fmt"
	"os"

	"github.com/lancewf/money-report-go/config"
	"github.com/lancewf/money-report-go/server"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Money-Report"
	app.Usage = "Money Report API Service"
	app.Copyright = "None"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		// Start Command
		{
			Name:   "start",
			Usage:  "Starts the Money Report API service",
			Action: startServer,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "es-url",
					Value: "http://elasticsearch:9200/",
					Usage: "Url to ElasticSearch (<protocol>://domain:<port>)/",
				},
				cli.IntFlag{
					Name:  "port",
					Value: 1234,
					Usage: "Port where the service will be listening",
				},
				cli.StringFlag{
					Name:  "host",
					Value: "0.0.0.0",
					Usage: "The ipaddress the service will be bound",
				},
				cli.StringFlag{
					Name:  "old-url",
					Value: "http://example.com",
					Usage: "The old running money programs base URL",
				},
			},
		},
	}

	app.Run(os.Args)
}

func startServer(ctx *cli.Context) error {

	cfg := config.GetDefault()
	cfg.ElasticsearchURL = ctx.String("es-url")
	cfg.Port = ctx.Int("port")
	cfg.Host = ctx.String("host")
	cfg.OldURL = ctx.String("old-url")
	fmt.Printf("starting server %o\n", cfg)

	server.Start(cfg)

	return nil
}
