package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	now := time.Now().UTC()
	str := ToString(now)
	now2, err := FromString(str)
	assert.Nil(err)
	str2 := ToString(now2)
	assert.Equal(str, str2)

	// wrong format
	_, err = FromString("wrongformat")
	assert.NotNil(err)
}
