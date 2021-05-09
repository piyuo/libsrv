package google

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoogle(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	assert.NotNil(Regions)
}
