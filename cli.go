package main

import (
	"fmt"
	"os"

	"github.com/matt-snider/voipms/api"
	vcli "github.com/matt-snider/voipms/cli"
	"github.com/urfave/cli"
)

/**
 * CLI Declaration
 */
func main() {
	app := cli.NewApp()
	app.Name = "voipms"
	app.Commands = []cli.Command{
		{
			Name:  "sms",
			Usage: "Manage SMS messages",
			Subcommands: []cli.Command{
				{
					Name:    "fetch",
					Aliases: []string{"get"},
					Usage:   "Fetch SMS messages",
					Action:  executeAction(vcli.FetchSms),
					Flags: []cli.Flag{
						cli.StringFlag{Name: "id"},
						cli.StringFlag{Name: "did"},
						cli.StringFlag{Name: "to"},
						cli.IntFlag{Name: "limit"},
						cli.StringFlag{
							Name:  "from-date",
							Usage: "Include message from this date (e.g 2017-01-31)",
						},
						cli.StringFlag{
							Name:  "to-date",
							Usage: "Include message up until this date (e.g 2017-01-31)",
						},
					},
				},

				{
					Name:      "send",
					Action:    executeAction(vcli.SendSms),
					ArgsUsage: "[did] [dest] [msg]",
					Usage:     "Send an SMS with a specific DID",
				},
			},
			// Default to fetch
			Action: executeAsDefault(vcli.FetchSms),
		},

		{
			Name:   "balance",
			Usage:  "Get balance info about this account (i.e. spending, call time, balance)",
			Action: executeAction(vcli.GetBalance),
		},

		{
			Name:   "ip",
			Usage:  "Get current IP address as seen by VoIP MS",
			Action: executeAction(vcli.GetIP),
		},
	}

	app.Run(os.Args)
}

/**
 * Helpers
 */
type command func(*api.VoipMsClient, *cli.Context) error

// Helper to provide client to command
func executeAction(fn command) cli.ActionFunc {
	username := os.Getenv("VOIPMS_USERNAME")
	password := os.Getenv("VOIPMS_PASSWORD")
	client := api.NewClient(username, password)
	return func(c *cli.Context) error {
		return fn(client, c)
	}
}

// Execute a command as the default command
// in a subcommand. Check the number of args
// passed to prevent an invalid subcommand
// triggering this. (e.g. voipms sms foo)
func executeAsDefault(fn command) cli.ActionFunc {
	username := os.Getenv("VOIPMS_USERNAME")
	password := os.Getenv("VOIPMS_PASSWORD")
	client := api.NewClient(username, password)
	return func(c *cli.Context) error {
		if c.Args().Present() {
			err := fmt.Sprintf("No command named '%s'", c.Args().First())
			return cli.NewExitError(err, -1)
		}
		return fn(client, c)
	}
}
