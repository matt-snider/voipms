package cli

import (
	"fmt"

	"github.com/matt-snider/voipms/api"
	"github.com/urfave/cli"
)

func GetIp(client *api.VoipMsClient, c *cli.Context) error {
	resp, err := client.GetIp()
	if err != nil {
		return cli.NewExitError(err, -1)
	}
	fmt.Printf("IP: %s\n", resp.IP)
	return nil
}
