package identifier

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUUID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	id := UUID()
	assert.NotEmpty(id)
}

func TestGoogleUUIDToString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	id := uuid.Must(uuid.NewRandom())
	token := GoogleUUIDToString(id)
	assert.NotEmpty(token)
}
