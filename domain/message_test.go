package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	t.Run("Validation fail", func(t *testing.T) {
		expectedErrorPayload := `{"body":"cannot be blank"}`
		m := Message{}
		err := m.Validate()
		assert.NotNil(t, err)
		errorPayload, err := json.Marshal(err)
		assert.Nil(t, err)
		assert.Equal(t, expectedErrorPayload, string(errorPayload))
	})

	t.Run("Validation ok", func(t *testing.T) {
		m := Message{Body: `{"type": "message"}`}
		err := m.Validate()
		assert.Nil(t, err)
	})

	t.Run("Enqueue", func(t *testing.T) {
		queue := Queue{
			ID:                      "my-queue",
			AckDeadlineSeconds:      60,
			MessageRetentionSeconds: 3600,
			DeliveryDelaySeconds:    10,
		}
		m := Message{Body: `{"type": "message"}`}
		now := time.Now().UTC()

		m.Enqueue(&queue, now)

		assert.NotEmpty(t, m.ID)
		assert.Equal(t, queue.ID, m.QueueID)
		assert.Equal(t, uint(0), m.DeliveryAttempts)
		assert.Equal(t, now.Add(time.Duration(queue.MessageRetentionSeconds)*time.Second), m.ExpiredAt)
		assert.Equal(t, now.Add(time.Duration(queue.DeliveryDelaySeconds)*time.Second), m.ScheduledAt)
		assert.Equal(t, now, m.CreatedAt)
		assert.Equal(t, now, m.UpdatedAt)
	})

	t.Run("DeliverySetup", func(t *testing.T) {
		queue := Queue{
			ID:                      "my-queue",
			AckDeadlineSeconds:      60,
			MessageRetentionSeconds: 3600,
			DeliveryDelaySeconds:    0,
		}
		m := Message{Body: `{"type": "message"}`}

		m.Enqueue(&queue, time.Now().UTC())
		now := time.Now().UTC()
		m.DeliverySetup(&queue, now)

		assert.Equal(t, uint(1), m.DeliveryAttempts)
		assert.Equal(t, now.Add(time.Duration(queue.AckDeadlineSeconds)*time.Second), m.ScheduledAt)
		assert.Equal(t, now, m.UpdatedAt)
	})

	t.Run("Ack", func(t *testing.T) {
		queue := Queue{
			ID:                      "my-queue",
			AckDeadlineSeconds:      60,
			MessageRetentionSeconds: 3600,
			DeliveryDelaySeconds:    0,
		}
		m := Message{Body: `{"type": "message"}`}

		m.Enqueue(&queue, time.Now().UTC())
		m.DeliverySetup(&queue, time.Now().UTC())
		now := time.Now().UTC()
		m.Ack(now)

		assert.Equal(t, now, m.ExpiredAt)
		assert.Equal(t, now, m.UpdatedAt)
	})

	t.Run("Nack", func(t *testing.T) {
		queue := Queue{
			ID:                      "my-queue",
			AckDeadlineSeconds:      60,
			MessageRetentionSeconds: 3600,
			DeliveryDelaySeconds:    0,
		}
		m := Message{Body: `{"type": "message"}`}

		m.Enqueue(&queue, time.Now().UTC())
		m.DeliverySetup(&queue, time.Now().UTC())
		now := time.Now().UTC()
		m.Nack(now, 100)

		assert.Equal(t, now.Add(time.Duration(100)*time.Second), m.ScheduledAt)
		assert.Equal(t, now, m.UpdatedAt)
	})
}
