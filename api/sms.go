package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Sms struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	SmsType string `json:"type"`
	DID     string `json:"did"`
	Contact string `json:"contact"`
	Message string `json:"message"`
}

type SmsResponse struct {
	VoipMsResponse
	SmsList []Sms `json:"sms"`
}

func (c *VoipMsClient) GetSms() (*SmsResponse, error) {
	resp, err := c.do("GET", "getSMS", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data *SmsResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	if err := toError(data.Status); err != nil {
		return nil, err
	}
	return data, nil
}

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
