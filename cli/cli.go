package main

import (
	"fmt"
	"os"

	"github.com/matt-snider/voipms/api"
)

func main() {
	username := os.Getenv("VOIPMS_USERNAME")
	password := os.Getenv("VOIPMS_PASSWORD")
	client := api.NewClient(username, password)
	resp, err := client.GetSms()
	if err != nil {
		fmt.Println("Error - ", err)
		os.Exit(-1)
	} else {
		fmt.Println("Response - ", resp)
	}
}
