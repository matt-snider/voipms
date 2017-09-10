package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Sms represents a message processed by voip.ms
type Sms struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	SmsType string `json:"type"`
	DID     string `json:"did"`
	Contact string `json:"contact"`
	Message string `json:"message"`
}

// SmsFilter describes which Sms messages to return
type SmsFilter struct {
	ID       string
	DID      string
	Contact  string
	FromDate time.Time
	ToDate   time.Time
	Limit    int
}

type smsResponse struct {
	VoipMsResponse
	SmsList []Sms `json:"sms"`
}

// GetSms returns an array of Sms instances.
//
// Currently this is just messages for the current day.
func (c *VoipMsClient) GetSms(filter SmsFilter) ([]Sms, error) {
	params := make(map[string]string)
	if filter.ID != "" {
		params["sms"] = filter.ID
	}
	if filter.DID != "" {
		params["did"] = filter.DID
	}
	if filter.Limit != 0 {
		params["limit"] = strconv.Itoa(filter.Limit)
	}
	if filter.Contact != "" {
		params["contact"] = filter.Contact
	}
	if filter.FromDate != (time.Time{}) {
		params["from"] = filter.FromDate.Format("2006-01-02")
	}
	if filter.ToDate != (time.Time{}) {
		params["to"] = filter.ToDate.Format("2006-01-02")
	}

	resp, err := c.do("GET", "getSMS", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data *smsResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	if err := toError(data.Status); err != nil {
		return nil, err
	}
	return data.SmsList, nil
}

// SendSms sends an SMS Message to the given dest from the provided DID.
func (c *VoipMsClient) SendSms(did, dest, msg string) error {
	resp, err := c.do(
		"POST",
		"sendSMS",
		map[string]string{
			"did":     url.QueryEscape(did),
			"dst":     url.QueryEscape(dest),
			"message": url.QueryEscape(msg),
		},
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data *VoipMsResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		return err
	}
	fmt.Println(data)
	if err := toError(data.Status); err != nil {
		return err
	}
	return nil
}
