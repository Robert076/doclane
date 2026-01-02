package requests

type GetUsersRequest struct {
	Limit   *int    `json:"limit"`
	Offset  *int    `json:"offset"`
	OrderBy *string `json:"order_by"`
	Order   *string `json:"order"`
}
