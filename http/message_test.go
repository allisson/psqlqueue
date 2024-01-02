package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/allisson/psqlqueue/domain"
)

func nilString() *string {
	var s *string = nil
	return s
}

func TestMessageHandler(t *testing.T) {
	t.Run("Create with invalid request", func(t *testing.T) {
		expectedPayload := `{"code":2,"message":"malformed request body"}`
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/queues/my-queue/messages", bytes.NewBuffer([]byte(`{`)))

		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusBadRequest, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Create with validation error", func(t *testing.T) {
		expectedPayload := `{"code":3,"message":"request validation failed","details":"body: cannot be blank."}`
		message := domain.Message{QueueID: "my-queue"}
		jsonMessage, _ := json.Marshal(&message)
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/queues/my-queue/messages", bytes.NewBuffer(jsonMessage))

		tc.messageService.On("Create", mock.Anything, &message).Return(message.Validate())
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusBadRequest, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Create", func(t *testing.T) {
		message := domain.Message{QueueID: "my-queue", Body: `{"message": true}`}
		jsonMessage, _ := json.Marshal(&message)
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/queues/my-queue/messages", bytes.NewBuffer(jsonMessage))

		tc.messageService.On("Create", mock.Anything, &message).Return(nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusNoContent, reqRec.Code)
	})

	t.Run("List", func(t *testing.T) {
		expectedPayload := `{"data":[{"id":"","queue_id":"my-queue","label":null,"body":"{\"message\": true}","attributes":null,"delivery_attempts":0,"created_at":"0001-01-01T00:00:00Z"},{"id":"","queue_id":"my-queue","label":null,"body":"{\"message\": true}","attributes":null,"delivery_attempts":0,"created_at":"0001-01-01T00:00:00Z"}],"limit":10}`
		message1 := domain.Message{QueueID: "my-queue", Body: `{"message": true}`}
		message2 := domain.Message{QueueID: "my-queue", Body: `{"message": true}`}
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/queues/my-queue/messages", nil)

		tc.messageService.On("List", mock.Anything, "my-queue", nilString(), uint(10)).Return([]*domain.Message{&message1, &message2}, nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusOK, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Ack", func(t *testing.T) {
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/v1/queues/my-queue/messages/message-id/ack", nil)

		tc.messageService.On("Ack", mock.Anything, "message-id").Return(nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusNoContent, reqRec.Code)
	})

	t.Run("Nack", func(t *testing.T) {
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/v1/queues/my-queue/messages/message-id/nack", bytes.NewBuffer([]byte(`{"visibility_timeout_seconds": 0}`)))

		tc.messageService.On("Nack", mock.Anything, "message-id", uint(0)).Return(nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusNoContent, reqRec.Code)
	})
}
