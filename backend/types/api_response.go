package types

type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Msg     string `json:"message,omitempty"`
	Err     string `json:"error,omitempty"`
}
