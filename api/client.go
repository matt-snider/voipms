package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type VoipMsClient struct {
	client   *http.Client
	apiUrl   string
	username string
	password string
}

type VoipMsResponse struct {
	Status string `json:"status"`
}

type IPResponse struct {
	VoipMsResponse
	IP string `json:"ip"`
}

const (
	userAgent = "go/voipms-client (https://github.com/matt-snider/voipms)"
)

func NewClient(username, password string) *VoipMsClient {
	return &VoipMsClient{
		client:   http.DefaultClient,
		apiUrl:   "https://voip.ms/api/v1/rest.php",
		username: username,
		password: password,
	}
}

func (c *VoipMsClient) newRequest(httpMethod, apiMethod string, params *map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(httpMethod, c.apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Build request depending on httpMethod
	form := url.Values{}
	if params != nil {
		for key, value := range *params {
			form.Add(key, value)
		}
	}
	form.Add("api_username", c.username)
	form.Add("api_password", c.password)
	form.Add("method", apiMethod)
	if httpMethod == http.MethodGet {
		req.URL.RawQuery = form.Encode()
	} else if httpMethod == http.MethodPost {
		req.PostForm = form
	} else {
		return nil, fmt.Errorf("Invalid http method for voipms api: %s", httpMethod)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)
	return req, nil
}

func (c *VoipMsClient) do(request *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *VoipMsClient) GetIp() (*IPResponse, error) {
	req, err := c.newRequest("GET", "getIP", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data *IPResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	if err := toError(data.Status); err != nil {
		return nil, err
	}
	return data, nil
}

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
	req, err := c.newRequest("GET", "getSMS", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
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
	params := &map[string]string{
		"did":     url.QueryEscape(did),
		"dst":     url.QueryEscape(dest),
		"message": url.QueryEscape(msg),
	}
	req, err := c.newRequest("POST", "sendSMS", params)
	fmt.Println(req.URL.RawQuery)
	if err != nil {
		return err
	}
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data *VoipMsResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		return err
	}
	if err := toError(data.Status); err != nil {
		return err
	}
	return nil
}

func toError(status string) error {
	var detail string

	switch status {
	case "success":
		detail = ""
	case "ip_not_enabled":
		detail = "This IP is not enabled for API use"
	case "invalid_method":
		detail = "This is not a valid Method"
	case "invalid_dst":
		detail = "This is not a valid Destination Number"
	case "missing_method":
		detail = "Method must be provided when using the REST/JSON API"
	case "missing_credentials":
		detail = "Username or Password was not provided"
	case "no_sms":
		detail = "There are no SMS messages"
	default:
		detail = status
	}

	if detail == "" {
		return nil
	}
	return fmt.Errorf("voipms: %s", detail)
}
