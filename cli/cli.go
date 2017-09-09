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
					fmt.Println("Error - ", err)
					os.Exit(-1)
				}

				fmt.Println("Resp %s", resp)
				return nil
			},
		},
	}
	app.Run(os.Args)
}
