package models

type TokenResponse struct {
	OngoingTokens []*Token         `json:"ongoingTokens"`
	WonTokens     []*WonToken      `json:"wonTokens"`
	ExpiredTokens []*TokenArchives `json:"expiredTokens"`
}

type APIResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
