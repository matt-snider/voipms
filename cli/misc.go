package cli

import (
	"fmt"

	"github.com/matt-snider/voipms/api"
	"github.com/urfave/cli"
)

func GetIP(client *api.VoipMsClient, c *cli.Context) error {
	ip, err := client.GetIP()
	if err != nil {
		return cli.NewExitError(err, -1)
	}
	fmt.Printf("IP: %s\n", ip)
	return nil
}
