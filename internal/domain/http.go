package domain

type HTTPError struct {
	Error string `json:"error"`
}

type HTTPok struct {
	Status string `json:"status"`
}
