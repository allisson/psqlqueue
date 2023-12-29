package http

type listResponse struct {
	Data   interface{} `json:"data"`
	Offset uint        `json:"offset,omitempty"`
	Limit  uint        `json:"limit"`
} //@name ListResponse
