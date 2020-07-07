package util

import (
	"encoding/binary"
	"math/rand"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
)

// UUID generates base58,max 22 character long,concise, unambiguous, URL-safe UUIDs string
//	return empty if something wrong
//
//	id, err := UUID() //PMty86Lju4PiaUAhspHYAn
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
//	id := SerialID16(uint16(42)) //4Go
//
func SerialID16(i uint16) string {
	i++ // avoid zero
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, uint16(i))
	return base58.Encode(bytes)
}

// SerialID32 use uint32 value between 1 to 4,294,967,295 to base58 5~6 character long,concise, unambiguous, URL-safe string
//
//	id := SerialID32(uint32(42)) //26kU7q
//
func SerialID32(i uint32) string {
	i++ // avoid zero
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(i))
	return base58.Encode(bytes)
}

// SerialID64 use uint64 value between 1 to 18,446,744,073,709,551,615 to base58 10~11 character long,concise, unambiguous, URL-safe string
//
//	id := SerialID64(uint64(42)) //8C9vbiDD9WF
//
func SerialID64(i uint64) string {
	i++ // avoid zero
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(i))
	return base58.Encode(bytes)
}

var randSrc = rand.NewSource(time.Now().UnixNano())

const letterBytes = "1234567890"
const (
	letterIdxBits = 4                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// OrderNumber format like credit card number, e.g. 0623-8476-6612-5355 , first 4 digit is current date ,rest is random number
//
//	id := OrderNumber() //0624-9128-0038-1148
//
func OrderNumber() string {
	n := 12
	sb := strings.Builder{}
	sb.Grow(19)
	sb.WriteString(time.Now().Format("0102"))
	sb.WriteString("-")
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		l := sb.Len()
		if l == 9 || l == 14 {
			sb.WriteString("-")
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}
