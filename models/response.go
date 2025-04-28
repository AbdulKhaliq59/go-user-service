package models

import "time"

type TokenResponse struct {
	OngoingTokens []*Token         `json:"ongoingTokens"`
	WonTokens     []*WonToken      `json:"wonTokens"`
	ExpiredTokens []*TokenArchives `json:"expiredTokens"`
}

type APIResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type Response struct {
	Timestamp time.Time   `json:"timestamp"`
	Message   string      `json:"message"`
	Status    int         `json:"status"`
	Data      interface{} `json:"data"`
}

type LoginResponseData struct {
	Token string `json:"token"`
}
