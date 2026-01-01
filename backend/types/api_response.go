package types

type APIResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"message,omitempty"`
	Err     string `json:"error,omitempty"`
	Token   string `json:"token,omitempty"`
}
