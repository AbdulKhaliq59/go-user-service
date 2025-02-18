package models

type TokenStats struct {
	Product     *Product `json:"product"`
	Draw        *Draw    `json:"draw"`
	NbrOfTokens int      `json:"nbrOfTokens"`
}
