package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscription(t *testing.T) {
	t.Run("Validation fail", func(t *testing.T) {
		expectedErrorPayload := `{"id":"must be in a valid format","queue_id":"must be in a valid format","topic_id":"must be in a valid format"}`
		subs := Subscription{ID: "my@invalid@id", TopicID: "my@invalid@id", QueueID: "my@invalid@id"}
		err := subs.Validate()
		assert.NotNil(t, err)
		errorPayload, err := json.Marshal(err)
		assert.Nil(t, err)
		assert.Equal(t, expectedErrorPayload, string(errorPayload))
	})

	t.Run("Validation ok", func(t *testing.T) {
		subs := Subscription{ID: "my-subscription", TopicID: "my-topic", QueueID: "my-queue"}
		err := subs.Validate()
		assert.Nil(t, err)
	})

	t.Run("ShouldCreateMessage", func(t *testing.T) {
		tests := []struct {
			subscription Subscription
			message      Message
			expected     bool
		}{
			{
				subscription: Subscription{},
				message:      Message{},
				expected:     true,
			},
			{
				subscription: Subscription{},
				message:      Message{Attributes: (map[string]string{"type": "message"})},
				expected:     true,
			},
			{
				subscription: Subscription{MessageFilters: map[string][]string{"type": {"message"}}},
				message:      Message{},
				expected:     false,
			},
			{
				subscription: Subscription{MessageFilters: map[string][]string{"type": {"message"}}},
				message:      Message{Attributes: map[string]string{"type": "message2"}},
				expected:     false,
			},
			{
				subscription: Subscription{MessageFilters: map[string][]string{"type": {"message"}}},
				message:      Message{Attributes: map[string]string{"type": "message"}},
				expected:     true,
			},
			{
				subscription: Subscription{MessageFilters: map[string][]string{"type": {"message", "message2"}}},
				message:      Message{Attributes: map[string]string{"type": "message"}},
				expected:     true,
			},
			{
				subscription: Subscription{MessageFilters: map[string][]string{"type": {"message", "message2"}}},
				message:      Message{Attributes: map[string]string{"type": "message2"}},
				expected:     true,
			},
			{
				subscription: Subscription{MessageFilters: map[string][]string{"type": {"message"}, "subtype": {"post"}}},
				message:      Message{Attributes: map[string]string{"type": "message"}},
				expected:     false,
			},
			{
				subscription: Subscription{MessageFilters: map[string][]string{"type": {"message"}, "subtype": {"post"}}},
				message:      Message{Attributes: map[string]string{"type": "message", "subtype": "comment"}},
				expected:     false,
			},
			{
				subscription: Subscription{MessageFilters: map[string][]string{"type": {"message"}, "subtype": {"post"}}},
				message:      Message{Attributes: map[string]string{"type": "message", "subtype": "post"}},
				expected:     true,
			},
			{
				subscription: Subscription{MessageFilters: map[string][]string{"type": {"message"}, "subtype": {"post", "comment"}}},
				message:      Message{Attributes: map[string]string{"type": "message", "subtype": "comment"}},
				expected:     true,
			},
		}

		for i := range tests {
			t.Run("", func(t *testing.T) {
				assert.Equal(t, tests[i].expected, tests[i].subscription.ShouldCreateMessage(&tests[i].message))
			})
		}
	})
}
