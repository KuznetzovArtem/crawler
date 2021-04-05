package domain

import "encoding/json"

type Request struct {
	Action       string          `json:"action"`
	ActionFields json.RawMessage `json:"fields"`
}
