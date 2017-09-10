package cli

import (
	"fmt"

	"github.com/matt-snider/voipms/api"
	"github.com/urfave/cli"
)

func FetchSms(client *api.VoipMsClient, c *cli.Context) error {
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
}

func SendSms(client *api.VoipMsClient, c *cli.Context) error {
	did := c.Args().Get(0)
	dest := c.Args().Get(1)
	msg := c.Args().Get(2)
	if did == "" {
		return cli.NewExitError("[did] not set", -1)
	} else if dest == "" {
		return cli.NewExitError("[dest] not set", -1)
	} else if msg == "" {
		return cli.NewExitError("[msg] not set", -1)
	}
	err := client.SendSms(did, dest, msg)
	if err != nil {
		return cli.NewExitError(err, -1)
	}
	fmt.Println("Message sent.")
	return nil
}
