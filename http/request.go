package http

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type listRequest struct {
	Offset uint `form:"offset"`
	Limit  uint `form:"limit"`
}

func newListRequestFromGIN(c *gin.Context) listRequest {
	request := listRequest{Offset: 0, Limit: 1}
	if err := c.ShouldBindQuery(&request); err != nil {
		slog.Warn("list request error", "error", err)
	}
	return request
}
