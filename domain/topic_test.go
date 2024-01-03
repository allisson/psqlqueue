package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopic(t *testing.T) {
	t.Run("Validation fail", func(t *testing.T) {
		expectedErrorPayload := `{"id":"must be in a valid format"}`
		topic := Topic{ID: "my@invalid@id"}
		err := topic.Validate()
		assert.NotNil(t, err)
		errorPayload, err := json.Marshal(err)
		assert.Nil(t, err)
		assert.Equal(t, expectedErrorPayload, string(errorPayload))
	})

	t.Run("Validation ok", func(t *testing.T) {
		topic := Topic{ID: "my-topic"}
		err := topic.Validate()
		assert.Nil(t, err)
	})
}
