package api

import "encoding/json"

type ipResponse struct {
	VoipMsResponse
	IP string `json:"ip"`
}

// GetIP returns the IPv4 address as seen by voip.ms
func (c *VoipMsClient) GetIP() (string, error) {
	resp, err := c.do("GET", "getIP", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data *ipResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		return "", err
	}
	if err := toError(data.Status); err != nil {
		return "", err
	}
	return data.IP, nil
}
