package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
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

func (c *VoipMsClient) do(httpMethod, apiMethod string, params map[string]string) (*http.Response, error) {
	req, err := c.newRequest(httpMethod, apiMethod, params)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *VoipMsClient) newRequest(httpMethod, apiMethod string, params map[string]string) (*http.Request, error) {
	if params == nil {
		params = make(map[string]string)
	}
	params["api_username"] = c.username
	params["api_password"] = c.password
	params["method"] = apiMethod

	var url string
	var contentType string
	var body = &bytes.Buffer{}
	if httpMethod == http.MethodGet {
		url = c.apiUrl + "?" + toUrlValues(params).Encode()
	} else if httpMethod == http.MethodPost {
		url = c.apiUrl
		writer := multipart.NewWriter(body)
		for k, v := range params {
			writer.WriteField(k, v)
		}
		contentType = writer.FormDataContentType()
		writer.Close()
	} else {
		return nil, fmt.Errorf("Invalid http method for voipms api: %s", httpMethod)
	}

	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)
	return req, nil
}

func toUrlValues(params map[string]string) url.Values {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return values
}

func (c *VoipMsClient) GetIp() (*IPResponse, error) {
	resp, err := c.do("GET", "getIP", nil)
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
	case "invalid_did":
		detail = "This is not a valid DID"
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
