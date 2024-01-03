package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/allisson/psqlqueue/domain"
)

func TestHealthCheckHandler(t *testing.T) {
	t.Run("Check", func(t *testing.T) {
		expectedPayload := `{"success":true}`
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/healthz", nil)

		tc.healthCheckService.On("Check", mock.Anything).Return(&domain.HealthCheck{Success: true}, nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusOK, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})
}
