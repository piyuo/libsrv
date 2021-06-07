package identifier

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
)

// UUID generates base58,max 22 character long,concise, unambiguous, URL-safe UUIDs string
//	return empty if something wrong
//
//	id := UUID() // PMty86Lju4PiaUAhspHYAn
//
func UUID() string {
	id := uuid.Must(uuid.NewRandom())
	return GoogleUUIDToString(id)
}

// GoogleUUIDToString convert google uuid to base58 string
//	return empty if something wrong
//
//	token := GoogleUUIDToString(id) // PMty86Lju4PiaUAhspHYAn
//
func GoogleUUIDToString(id uuid.UUID) string {
	bytes, err := id.MarshalBinary()
	if err != nil {
		return ""
	}
	return base58.Encode(bytes)
}

// GoogleUUIDFromString convert string to google uuid
//
//	id, err := GoogleUUIDToString(id) // PMty86Lju4PiaUAhspHYAn
//
func GoogleUUIDFromString(token string) (uuid.UUID, error) {
	bytes := base58.Decode(token)
	var id uuid.UUID
	if err := id.UnmarshalBinary(bytes); err != nil {
		return id, err
	}
	return id, nil
}
