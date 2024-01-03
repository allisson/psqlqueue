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

func TestSubscriptionHandler(t *testing.T) {
	t.Run("Create with invalid request", func(t *testing.T) {
		expectedPayload := `{"code":2,"message":"malformed request body"}`
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/subscriptions", bytes.NewBuffer([]byte(`{`)))

		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusBadRequest, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Create with validation error", func(t *testing.T) {
		expectedPayload := `{"code":3,"message":"request validation failed","details":"id: must be in a valid format; queue_id: cannot be blank; topic_id: cannot be blank."}`
		subscription := domain.Subscription{ID: "my@subscription"}
		jsonSubscription, _ := json.Marshal(&subscription)
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/subscriptions", bytes.NewBuffer(jsonSubscription))

		tc.subscriptionService.On("Create", mock.Anything, &subscription).Return(subscription.Validate())
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusBadRequest, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Create", func(t *testing.T) {
		expectedPayload := `{"id":"my-subscription","topic_id":"my-topic","queue_id":"my-queue","message_filters":null,"created_at":"0001-01-01T00:00:00Z"}`
		subscription := domain.Subscription{ID: "my-subscription", TopicID: "my-topic", QueueID: "my-queue"}
		jsonSubscription, _ := json.Marshal(&subscription)
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/subscriptions", bytes.NewBuffer(jsonSubscription))

		tc.subscriptionService.On("Create", mock.Anything, &subscription).Return(nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusCreated, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Get with object not found", func(t *testing.T) {
		expectedPayload := `{"code":10,"message":"subscription not found"}`
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/subscriptions/my-subscription", nil)

		tc.subscriptionService.On("Get", mock.Anything, "my-subscription").Return(nil, domain.ErrSubscriptionNotFound)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusNotFound, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Get", func(t *testing.T) {
		expectedPayload := `{"id":"my-subscription","topic_id":"my-topic","queue_id":"my-queue","message_filters":null,"created_at":"0001-01-01T00:00:00Z"}`
		subscription := domain.Subscription{ID: "my-subscription", TopicID: "my-topic", QueueID: "my-queue"}
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/subscriptions/my-subscription", nil)

		tc.subscriptionService.On("Get", mock.Anything, subscription.ID).Return(&subscription, nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusOK, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("List", func(t *testing.T) {
		expectedPayload := `{"data":[{"id":"my-subscription-1","topic_id":"my-topic","queue_id":"my-queue-1","message_filters":null,"created_at":"0001-01-01T00:00:00Z"},{"id":"my-subscription-2","topic_id":"my-topic","queue_id":"my-queue-2","message_filters":null,"created_at":"0001-01-01T00:00:00Z"}],"limit":1}`
		subscription1 := domain.Subscription{ID: "my-subscription-1", TopicID: "my-topic", QueueID: "my-queue-1"}
		subscription2 := domain.Subscription{ID: "my-subscription-2", TopicID: "my-topic", QueueID: "my-queue-2"}
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/subscriptions", nil)

		tc.subscriptionService.On("List", mock.Anything, uint(0), uint(1)).Return([]*domain.Subscription{&subscription1, &subscription2}, nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusOK, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Delete with object not found", func(t *testing.T) {
		expectedPayload := `{"code":10,"message":"subscription not found"}`
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/v1/subscriptions/my-subscription", nil)

		tc.subscriptionService.On("Delete", mock.Anything, "my-subscription").Return(domain.ErrSubscriptionNotFound)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusNotFound, reqRec.Code)
		assert.Equal(t, expectedPayload, reqRec.Body.String())
	})

	t.Run("Delete", func(t *testing.T) {
		tc := makeTestContext(t)
		reqRec := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/v1/subscriptions/my-subscription", nil)

		tc.subscriptionService.On("Delete", mock.Anything, "my-subscription").Return(nil)
		tc.router.ServeHTTP(reqRec, req)

		assert.Equal(t, http.StatusNoContent, reqRec.Code)
	})
}
