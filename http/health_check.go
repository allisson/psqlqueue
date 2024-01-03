package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/allisson/psqlqueue/domain"
)

// nolint:unused
type healthCheckResponse struct {
	Success bool `json:"success"`
} //@name HealthCheckResponse

// HealthCheckHandler exposes a REST API for domain.HealthCheckService.
type HealthCheckHandler struct {
	healthCheckService domain.HealthCheckService
}

// Execute a health check.
//
//	@Summary	Execute a health check
//	@Tags		health-check
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	healthCheckResponse
//	@Failure	500	{object}	errorResponse
//	@Router		/healthz [get]
func (h *HealthCheckHandler) Check(c *gin.Context) {
	healthCheck, err := h.healthCheckService.Check(c.Request.Context())
	if err != nil {
		er := parseServiceError("healthCheckService", "Check", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.JSON(http.StatusOK, &healthCheck)
}

// NewHealthCheckHandler returns a new HealthCheckHandler.
func NewHealthCheckHandler(healthCheckService domain.HealthCheckService) *HealthCheckHandler {
	return &HealthCheckHandler{healthCheckService: healthCheckService}
}
