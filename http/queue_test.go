package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/allisson/psqlqueue/domain"
	"github.com/allisson/psqlqueue/mocks"
)

type testContext struct {
	queueService        *mocks.QueueService
	queueHandler        *QueueHandler
	messageService      *mocks.MessageService
	messageHandler      *MessageHandler
	topicService        *mocks.TopicService
	topicHandler        *TopicHandler
	subscriptionService *mocks.SubscriptionService
	subscriptionHandler *SubscriptionHandler
	router              *gin.Engine
}

func makeTestContext(t *testing.T) *testContext {
	logger := domain.NewLogger("debug")
	queueService := mocks.NewQueueService(t)
	queueHandler := NewQueueHandler(queueService)
	messageService := mocks.NewMessageService(t)
	messageHandler := NewMessageHandler(messageService)
	topicService := mocks.NewTopicService(t)
	topicHandler := NewTopicHandler(topicService)
	subscriptionService := mocks.NewSubscriptionService(t)
	subscriptionHandler := NewSubscriptionHandler(subscriptionService)
	router := SetupRouter(logger, queueHandler, messageHandler, topicHandler, subscriptionHandler)
	return &testContext{
		queueService:        queueService,
		queueHandler:        queueHandler,
		messageService:      messageService,
		messageHandler:      messageHandler,
		topicService:        topicService,
		topicHandler:        topicHandler,
		subscriptionService: subscriptionService,
		subscriptionHandler: subscriptionHandler,
		router:              router,
	}
}

func TestQueueHandler(t *testing.T) {
	t.Run("Create with invalid request", func(t *testing.T) {
		expectedPayload := `{"code":2,"message":"malformed request body"}`
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/queues", bytes.NewBuffer([]byte(`{`)))

		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusBadRequest, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Create with validation error", func(t *testing.T) {
		expectedPayload := `{"code":3,"message":"request validation failed","details":"ack_deadline_seconds: cannot be blank; id: must be in a valid format; message_retention_seconds: cannot be blank."}`
		queue := domain.Queue{ID: "my@queue"}
		jsonQueue, _ := json.Marshal(&queue)
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/queues", bytes.NewBuffer(jsonQueue))

		tc.queueService.On("Create", mock.Anything, &queue).Return(queue.Validate())
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusBadRequest, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Create", func(t *testing.T) {
		expectedPayload := `{"id":"my-queue","ack_deadline_seconds":0,"message_retention_seconds":0,"delivery_delay_seconds":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`
		queue := domain.Queue{ID: "my-queue"}
		jsonQueue, _ := json.Marshal(&queue)
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/queues", bytes.NewBuffer(jsonQueue))

		tc.queueService.On("Create", mock.Anything, &queue).Return(nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusCreated, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Get with object not found", func(t *testing.T) {
		expectedPayload := `{"code":5,"message":"queue not found"}`
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/queues/my-queue", nil)

		tc.queueService.On("Get", mock.Anything, "my-queue").Return(nil, domain.ErrQueueNotFound)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusNotFound, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Get", func(t *testing.T) {
		expectedPayload := `{"id":"my-queue","ack_deadline_seconds":0,"message_retention_seconds":0,"delivery_delay_seconds":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`
		queue := domain.Queue{ID: "my-queue"}
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/queues/my-queue", nil)

		tc.queueService.On("Get", mock.Anything, queue.ID).Return(&queue, nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusOK, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Update with invalid request", func(t *testing.T) {
		expectedPayload := `{"code":2,"message":"malformed request body"}`
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/v1/queues/my-queue", bytes.NewBuffer([]byte(`{`)))

		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusBadRequest, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Update with validation error", func(t *testing.T) {
		expectedPayload := `{"code":3,"message":"request validation failed","details":"ack_deadline_seconds: cannot be blank; message_retention_seconds: cannot be blank."}`
		queue := domain.Queue{ID: "my-queue"}
		jsonQueue, _ := json.Marshal(&queue)
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/v1/queues/my-queue", bytes.NewBuffer(jsonQueue))

		tc.queueService.On("Update", mock.Anything, &queue).Return(queue.Validate())
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusBadRequest, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Update", func(t *testing.T) {
		expectedPayload := `{"id":"my-queue","ack_deadline_seconds":0,"message_retention_seconds":0,"delivery_delay_seconds":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`
		queue := domain.Queue{ID: "my-queue"}
		jsonQueue, _ := json.Marshal(&queue)
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/v1/queues/my-queue", bytes.NewBuffer(jsonQueue))

		tc.queueService.On("Update", mock.Anything, &queue).Return(nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusOK, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("List", func(t *testing.T) {
		expectedPayload := `{"data":[{"id":"my-queue-1","ack_deadline_seconds":0,"message_retention_seconds":0,"delivery_delay_seconds":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},{"id":"my-queue-2","ack_deadline_seconds":0,"message_retention_seconds":0,"delivery_delay_seconds":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}],"limit":1}`
		queue1 := domain.Queue{ID: "my-queue-1"}
		queue2 := domain.Queue{ID: "my-queue-2"}
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/queues", nil)

		tc.queueService.On("List", mock.Anything, uint(0), uint(1)).Return([]*domain.Queue{&queue1, &queue2}, nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusOK, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Delete with object not found", func(t *testing.T) {
		expectedPayload := `{"code":5,"message":"queue not found"}`
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/v1/queues/my-queue", nil)

		tc.queueService.On("Delete", mock.Anything, "my-queue").Return(domain.ErrQueueNotFound)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusNotFound, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Delete", func(t *testing.T) {
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/v1/queues/my-queue", nil)

		tc.queueService.On("Delete", mock.Anything, "my-queue").Return(nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusNoContent, reqRec.Code)
	})

	t.Run("Stats", func(t *testing.T) {
		expectedPayload := `{"num_undelivered_messages":0,"oldest_unacked_message_age_seconds":0}`
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/queues/my-queue/stats", nil)

		tc.queueService.On("Stats", mock.Anything, "my-queue").Return(&domain.QueueStats{}, nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusOK, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Purge", func(t *testing.T) {
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/v1/queues/my-queue/purge", nil)

		tc.queueService.On("Purge", mock.Anything, "my-queue").Return(nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusNoContent, reqRec.Code)
	})

	t.Run("Cleanup", func(t *testing.T) {
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/v1/queues/my-queue/cleanup", nil)

		tc.queueService.On("Cleanup", mock.Anything, "my-queue").Return(nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusNoContent, reqRec.Code)
	})
}
