package data

type FailureResponse struct {
	Detail string `json:"detail"`
	Code   string `json:"code"`
	Status bool   `json:"status"`
}
