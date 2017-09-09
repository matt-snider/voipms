package main

import (
	"fmt"
	"os"

	"github.com/matt-snider/voipms/api"
	"github.com/urfave/cli"
)

func main() {
	username := os.Getenv("VOIPMS_USERNAME")
	password := os.Getenv("VOIPMS_PASSWORD")
	client := api.NewClient(username, password)

	// CLI App
	app := cli.NewApp()
	app.Name = "voipms"
	app.Commands = []cli.Command{
		{
			Name:  "getsms",
			Usage: "Get SMS's",
			Action: func(c *cli.Context) error {
				resp, err := client.GetSms()
				if err != nil {
					return cli.NewExitError(err, -1)
				}

				lineFmt := "%-9s %-20s %-8s %-10s %s\n"
				fmt.Printf(lineFmt, "ID", "Date", "Action", "Contact", "Message")
				for _, sms := range resp.SmsList {
					var action string
					if sms.SmsType == "0" {
						action = "sent"
					} else {
						action = "received"
					}

					fmt.Printf(lineFmt, sms.ID, sms.Date, action, sms.Contact, sms.Message)
				}
				return nil
			},
		},

		{
			Name:  "getip",
			Usage: "Get current IP address as seen by VoIP MS",
			Action: func(c *cli.Context) error {
				resp, err := client.GetIp()
				if err != nil {
					return cli.NewExitError(err, -1)
				}
				fmt.Printf("IP: %s\n", resp.IP)
				return nil
			},
		},
	}
	app.Run(os.Args)
}
