package entity

import "encoding/json"

type Wallet struct {
	UserId string
	Chips  int64 `json:"chips"`
}

func ParseWallet(payload string) (Wallet, error) {
	w := Wallet{}
	err := json.Unmarshal([]byte(payload), &w)
	return w, err
}
