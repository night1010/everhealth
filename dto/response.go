package dto

type Response struct {
	Data        any    `json:"data,omitempty"`
	Error       any    `json:"error,omitempty"`
	CurrentPage *int   `json:"current_page,omitempty"`
	CurrentItem *int   `json:"current_item,omitempty"`
	TotalPage   *int   `json:"total_page,omitempty"`
	TotalItem   *int   `json:"total_item,omitempty"`
	Message     string `json:"message,omitempty"`
}
