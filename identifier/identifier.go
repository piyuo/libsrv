package util

import (
	"encoding/binary"

	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
)

// UUID generates base58,max 22 character long,concise, unambiguous, URL-safe UUIDs string
//	return empty if something wrong
//
//	id, err := UUID() // PMty86Lju4PiaUAhspHYAn
//
func UUID() string {
	id := uuid.Must(uuid.NewRandom())
	bytes, err := id.MarshalBinary()
	if err != nil {
		return ""
	}
	return base58.Encode(bytes)
}

// SerialID16 use uint16 value between 1 to 65,535 to base58 10~11 character long,concise, unambiguous, URL-safe string
//
//	id := SerialID16(uint16(42)) // 4Go
//
func SerialID16(i uint16) string {
	i++ // avoid zero
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, i)
	return base58.Encode(bytes)
}

// SerialID32 use uint32 value between 1 to 4,294,967,295 to base58 5~6 character long,concise, unambiguous, URL-safe string
//
//	id := SerialID32(uint32(42)) // 26kU7q
//
func SerialID32(i uint32) string {
	i++ // avoid zero
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, i)
	return base58.Encode(bytes)
}

// SerialID64 use int64 value between 1 to 18,446,744,073,709,551,615 to base58 10~11 character long,concise, unambiguous, URL-safe string
//
//	id := SerialID64(uint64(42)) //8C9vbiDD9WF
//
func SerialID64(i uint64) string {
	i++ // avoid zero
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, i)
	return base58.Encode(bytes)
}
