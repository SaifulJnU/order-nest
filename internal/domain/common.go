package domain

type HttpResponse struct {
	Message string      `json:"message,omitempty"`
	Type    string      `json:"type,omitempty"`
	Code    int         `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
