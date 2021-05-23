package digit

import (
	"testing"

	"github.com/piyuo/libsrv/identifier"
	"github.com/stretchr/testify/assert"
)

func TestGob(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	sample := &Sample{}
	sample.Value = 43
	sample.Text = "gob-" + identifier.RandomNumber(90)
	encoded, err := Encode(sample)
	assert.Nil(err)
	assert.NotNil(encoded)

	sample2 := &Sample{}
	err = Decode(encoded, sample2)
	assert.Nil(err)
	assert.NotNil(sample2)

	assert.Equal(sample.Text, sample2.Text)
	assert.Equal(sample.Value, sample2.Value)
}

func TestGobNil(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	encoded, err := Encode(nil)
	assert.Nil(err)
	assert.Nil(encoded)

	sample2 := &Sample{}
	err = Decode(nil, sample2)
	assert.Nil(err)

	err = Decode(nil, nil)
	assert.NotNil(err)
}

type Sample struct {
	Value int
	Text  string
}

func BenchmarkEncode(b *testing.B) {
	sample := &Sample{}
	sample.Value = 43
	sample.Text = "gob-" + identifier.RandomNumber(90)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(sample)
	}
}

func BenchmarkDecode(b *testing.B) {
	sample := &Sample{}
	sample.Value = 43
	sample.Text = "gob-" + identifier.RandomNumber(90)
	encoded, _ := Encode(sample)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(encoded, sample)
	}
}
