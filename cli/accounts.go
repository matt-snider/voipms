package cli

import (
	"fmt"
	"strconv"

	"github.com/matt-snider/voipms/api"
	"github.com/urfave/cli"
)

func GetBalance(client *api.VoipMsClient, c *cli.Context) error {
	// balanceInfo, err := client.GetBalance()
	// if err != nil {
	// 	return cli.NewExitError(err, -1)
	// }
	balanceInfo := api.BalanceInfo{
		Balance:    200,
		SpentTotal: 12.200,
		SpentToday: 0,
		CallsTotal: 10,
		CallsToday: 0,
		TimeTotal:  "10:00",
		TimeToday:  "0",
	}

	// Balance
	fmt.Fprintf(c.App.Writer, "Current balance: $%.2f\n\n", balanceInfo.Balance)

	// Stats
	spent := [...]string{
		strconv.FormatFloat(balanceInfo.SpentToday, 'f', 2, 64),
		strconv.FormatFloat(balanceInfo.SpentTotal, 'f', 2, 64),
		"Spent",
	}
	called := [...]string{
		strconv.Itoa(balanceInfo.CallsToday),
		strconv.Itoa(balanceInfo.CallsTotal),
		"Calls",
	}
	times := [...]string{
		string(balanceInfo.TimeToday),
		balanceInfo.TimeTotal,
		"Time",
	}

	// Format with enough room for all values in the column
	lineFmt := "%-7s " +
		"%-" + strconv.Itoa(getMaxLength(spent[:])+3) + "s " +
		"%-" + strconv.Itoa(getMaxLength(called[:])+3) + "s " +
		"%-" + strconv.Itoa(getMaxLength(times[:])+3) + "s\n"

	fmt.Fprintf(c.App.Writer, lineFmt, "", "Spent", "Calls", "Time")
	fmt.Fprintf(c.App.Writer, lineFmt, "Today", spent[0], called[0], times[0])
	fmt.Fprintf(c.App.Writer, lineFmt, "Total", spent[1], called[1], times[1])
	return nil
}

func getMaxLength(values []string) int {
	max := -1
	for _, value := range values {
		length := len(value)
		if length > max {
			max = length
		}
	}
	return max
}
