package domain

type Response struct {
	Error  string      `json:"error"`
	Result interface{} `json:"result"`
}
