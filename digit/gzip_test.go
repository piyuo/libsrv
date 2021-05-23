package digit

import (
	"testing"

	"github.com/piyuo/libsrv/identifier"
	"github.com/stretchr/testify/assert"
)

func TestGzip(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	str := "gzip-" + identifier.RandomNumber(90)
	bytes := []byte(str)

	zipped, err := Compress(bytes)
	assert.Nil(err)
	assert.NotNil(zipped)

	unzipped, err := Decompress(zipped)
	assert.Nil(err)
	assert.NotNil(unzipped)

	str2 := string(unzipped)
	assert.Equal(str, str2)
}

func TestGzipNullBytes(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	// no error compress nil or empty bytes
	zipped, err := Compress(nil)
	assert.Nil(err)
	assert.Nil(zipped)

	unzipped, err := Decompress(nil)
	assert.Nil(err)
	assert.Nil(unzipped)
}

func BenchmarkCompress(b *testing.B) {
	str := "benchmark-compress-" + identifier.RandomNumber(99)
	bytes := []byte(str)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Compress(bytes)
	}
}

func BenchmarkDecompress(b *testing.B) {
	str := "benchmark-decompress-" + identifier.RandomNumber(99)
	bytes := []byte(str)
	zipped, _ := Compress(bytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decompress(zipped)
	}
}
