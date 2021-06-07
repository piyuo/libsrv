package identifier

import (
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
