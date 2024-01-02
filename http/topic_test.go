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

func TestTopicHandler(t *testing.T) {
	t.Run("Create with invalid request", func(t *testing.T) {
		expectedPayload := `{"code":2,"message":"malformed request body"}`
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer([]byte(`{`)))

		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusBadRequest, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Create with validation error", func(t *testing.T) {
		expectedPayload := `{"code":3,"message":"request validation failed","details":"id: must be in a valid format."}`
		topic := domain.Topic{ID: "my@topic"}
		jsonTopic, _ := json.Marshal(&topic)
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(jsonTopic))

		tc.topicService.On("Create", mock.Anything, &topic).Return(topic.Validate())
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusBadRequest, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Create", func(t *testing.T) {
		expectedPayload := `{"id":"my-topic","created_at":"0001-01-01T00:00:00Z"}`
		topic := domain.Topic{ID: "my-topic"}
		jsonTopic, _ := json.Marshal(&topic)
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(jsonTopic))

		tc.topicService.On("Create", mock.Anything, &topic).Return(nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusCreated, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Get with object not found", func(t *testing.T) {
		expectedPayload := `{"code":8,"message":"topic not found"}`
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/topics/my-topic", nil)

		tc.topicService.On("Get", mock.Anything, "my-topic").Return(nil, domain.ErrTopicNotFound)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusNotFound, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Get", func(t *testing.T) {
		expectedPayload := `{"id":"my-topic","created_at":"0001-01-01T00:00:00Z"}`
		topic := domain.Topic{ID: "my-topic"}
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/topics/my-topic", nil)

		tc.topicService.On("Get", mock.Anything, topic.ID).Return(&topic, nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusOK, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("List", func(t *testing.T) {
		expectedPayload := `{"data":[{"id":"my-topic-1","created_at":"0001-01-01T00:00:00Z"},{"id":"my-topic-2","created_at":"0001-01-01T00:00:00Z"}],"limit":1}`
		topic1 := domain.Topic{ID: "my-topic-1"}
		topic2 := domain.Topic{ID: "my-topic-2"}
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/topics", nil)

		tc.topicService.On("List", mock.Anything, uint(0), uint(1)).Return([]*domain.Topic{&topic1, &topic2}, nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusOK, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Delete with object not found", func(t *testing.T) {
		expectedPayload := `{"code":8,"message":"topic not found"}`
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/v1/topics/my-topic", nil)

		tc.topicService.On("Delete", mock.Anything, "my-topic").Return(domain.ErrTopicNotFound)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusNotFound, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Delete", func(t *testing.T) {
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/v1/topics/my-topic", nil)

		tc.topicService.On("Delete", mock.Anything, "my-topic").Return(nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusNoContent, reqRec.Code)
	})

	t.Run("CreateMessage", func(t *testing.T) {
		message := domain.Message{QueueID: "my-queue", Body: `{"message": true}`}
		jsonMessage, _ := json.Marshal(&message)
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/topics/my-topic/messages", bytes.NewBuffer(jsonMessage))

		tc.topicService.On("CreateMessage", mock.Anything, "my-topic", &message).Return(nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusNoContent, reqRec.Code)
	})
}
