package data

type BaseSuccessResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

type PaginationResponse struct {
	CurrentPage  string `json:"current_page"`
	TotalPages   int    `json:"total_pages"`
	TotalItems   int    `json:"total_items"`
	HasNext      bool   `json:"has_next"`
	HashPrevious bool   `json:"has_previous"`
}

type BaseSuccessListResponse struct {
	Code    string             `json:"code"`
	Message string             `json:"message,omitempty"`
	Page    PaginationResponse `json:"page"`
	Data    interface{}        `json:"data"`
}
