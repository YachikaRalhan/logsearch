package model

type ErrorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Status  string `json:"status"`
}

type ShortnerResponse struct {
	LogLines []string `json:"logLines"`
}
