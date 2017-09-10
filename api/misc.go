package api

import "encoding/json"

type IPResponse struct {
	VoipMsResponse
	IP string `json:"ip"`
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
