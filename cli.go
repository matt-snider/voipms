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
					Action:  executeAction(vcli.FetchSms),
					Usage:   "Fetch SMS messages",
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
			Name:   "ip",
			Usage:  "Get current IP address as seen by VoIP MS",
			Action: executeAction(vcli.GetIp),
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
