package identifier

import (
	"encoding/binary"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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

var randSrc rand.Source

const numberBytes = "1234567890"
const (
	numberIdxBits = 4                    // 6 bits to represent a letter index
	numberIdxMask = 1<<numberIdxBits - 1 // All 1-bits, as many as letterIdxBits
	numberIdxMax  = 63 / numberIdxBits   // # of letter indices fitting in 63 bits
)

// RandomNumber return number string on given digit
//
//	id := RandomNumber(6) //062448
//
func RandomNumber(digit int) string {
	if randSrc == nil {
		randSrc = rand.NewSource(time.Now().UnixNano())
	}

	sb := strings.Builder{}
	sb.Grow(digit)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := digit-1, randSrc.Int63(), numberIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), numberIdxMax
		}
		if idx := int(cache & numberIdxMask); idx < len(numberBytes) {
			sb.WriteByte(numberBytes[idx])
			i--
		}
		cache >>= numberIdxBits
		remain--
	}
	return sb.String()
}

// NotIdenticalRandomNumber return number string that avoid identical
//
//	id := RandomNumber(6) //062448
//
func NotIdenticalRandomNumber(digit int) string {

	for i := 0; i < 10; i++ {
		str := RandomNumber(digit)
		if !IsNumberStringIdentical(str) {
			return str
		}
	}
	return RandomNumber(digit)
}

// IsNumberStringIdentical return true has only 2 digit different
//
//	 IsNumberStringIdentical("111111") //true
//	 IsNumberStringIdentical("111112") //true
//	 IsNumberStringIdentical("111124") //false
//
func IsNumberStringIdentical(str string) bool {
	diffCount := 0
	c := str[0]
	for i := 1; i < len(str); i++ {
		if str[i] != c {
			diffCount++
		}
		c = str[i]
	}
	return diffCount <= 1
}

// letterBytes use in RandomString
//
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandomString return random string on given digit
//
//	id := RandomString(2) //Ax
//
func RandomString(n int) string {
	if randSrc == nil {
		randSrc = rand.NewSource(time.Now().UnixNano())
	}

	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}

// MapID generate unique serial id in map
//
//	m := map[string] string{}
//	id, err := MapID() // "1"
//
func MapID(m map[string]string) (string, error) {
	if len(m) == 0 {
		return "1", nil
	}

	max := 0
	for k := range m {
		i, err := strconv.Atoi(k)
		if err != nil {
			return "", errors.Wrap(err, "key to number")
		}
		if i > max {
			max = i
		}
	}
	return strconv.Itoa(max + 1), nil
}
