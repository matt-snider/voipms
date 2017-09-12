package api

import "encoding/json"

// BalanceInfo contains info about the amount spent, calls made,
// and time in calls
type BalanceInfo struct {
	Balance    float64     `json:"current_balance,string"`
	SpentTotal float64     `json:"spent_total"`
	SpentToday float64     `json:"spent_today"`
	CallsTotal int         `json:"calls_total"`
	CallsToday int         `json:"calls_today"`
	TimeTotal  string      `json:"time_total"`
	TimeToday  json.Number `json:"time_today,Number"`
}

type balanceResponse struct {
	VoipMsResponse
	Balance BalanceInfo
}

// GetBalance gets information about the spending, calls made and time spent
// in calls. See BalanceInfo
func (c *VoipMsClient) GetBalance() (*BalanceInfo, error) {
	resp, err := c.do("GET", "getBalance", map[string]string{"advanced": "true"})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data *balanceResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	if err := toError(data.Status); err != nil {
		return nil, err
	}
	return &data.Balance, nil
}
