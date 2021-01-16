package region

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegion(t *testing.T) {
	assert := assert.New(t)
	assert.NotEmpty(Current)
}
