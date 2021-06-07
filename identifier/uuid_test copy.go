package identifier

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUUID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	id := UUID()
	assert.NotEmpty(id)
}
