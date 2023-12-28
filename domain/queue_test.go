package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	t.Run("Validation fail", func(t *testing.T) {
		expectedErrorPayload := `{"ack_deadline_seconds":"cannot be blank","id":"must be in a valid format","message_retention_seconds":"cannot be blank"}`
		queue := Queue{
			ID:                      "my@invalid@id",
			AckDeadlineSeconds:      0,
			MessageRetentionSeconds: 0,
			DeliveryDelaySeconds:    0,
		}
		err := queue.Validate()
		assert.NotNil(t, err)
		errorPayload, err := json.Marshal(err)
		assert.Nil(t, err)
		assert.Equal(t, expectedErrorPayload, string(errorPayload))
	})

	t.Run("Validation ok", func(t *testing.T) {
		queue := Queue{
			ID:                      "my-queue",
			AckDeadlineSeconds:      60,
			MessageRetentionSeconds: 3600,
			DeliveryDelaySeconds:    0,
		}
		err := queue.Validate()
		assert.Nil(t, err)
	})
}
