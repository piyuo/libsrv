package digit

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"math"

	"github.com/pkg/errors"
)

// compress        8510     131032 ns/op    814861 B/op       24 allocs/op,compress
// decompress    138428       8525 ns/op     41744 B/op        7 allocs/op

// Compress byte array, no error if compress nil or empty byte array
//
//	zipped, err := Compress(bytes)
//
func Compress(originBytes []byte) ([]byte, error) {
	if len(originBytes) == 0 {
		return originBytes, nil
	}

	zipBuf := bytes.Buffer{}
	zipped := gzip.NewWriter(&zipBuf)
	_, err := zipped.Write(originBytes)
	if err != nil {
		return nil, errors.Wrap(err, "compress")
	}
	if err := zipped.Close(); err != nil {
		return nil, errors.Wrap(err, "close gzip")
	}
	zipBytes := zipBuf.Bytes()
	ratio := math.Round((1 - (float64(len(zipBytes)) / float64(len(originBytes)))) * float64(100))
	warning := "gzip"
	if ratio < 0 {
		warning = "warning:" + warning
	}
	fmt.Printf("%s %v bytes to %v, ratio %v%%\n", warning, len(originBytes), len(zipBytes), ratio)
	return zipBytes, nil
}

// Decompress byte array, no error if decompress nil or empty byte array
//
//	unzipped, err := Decompress(zipped)
//
func Decompress(zipped []byte) ([]byte, error) {
	if len(zipped) == 0 {
		return zipped, nil
	}

	rdr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, errors.Wrap(err, "new gzip")
	}

	data, err := ioutil.ReadAll(rdr)
	if err != nil {
		return nil, errors.Wrap(err, "decompress")
	}
	if err := rdr.Close(); err != nil {
		return nil, errors.Wrap(err, "close gzip")
	}
	return data, nil
}
