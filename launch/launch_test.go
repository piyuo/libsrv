package launch

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	assert := assert.New(t)
	assert.NotPanics(Checklist)

	os.Setenv("NAME", "")
	assert.Panics(Checklist)
	os.Setenv("NAME", "not empty")

	os.Setenv("REGION", "")
	assert.Panics(Checklist)
	os.Setenv("REGION", "not empty")

	os.Setenv("BRANCH", "")
	assert.Panics(Checklist)
	os.Setenv("BRANCH", "not empty")
}
