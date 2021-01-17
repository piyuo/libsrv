package identifier

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrderNumber(t *testing.T) {
	assert := assert.New(t)
	rand.Seed(time.Now().UnixNano())
	num := OrderNumber()
	str := OrderNumberToString(num)
	valid := OrderNumberIsValid(str)
	assert.True(valid)
	assert.NotEmpty(str)
	retNum, err := OrderNumberFromString(str)
	assert.Nil(err)
	assert.Equal(num, retNum)
	//fmt.Printf("%v\n", str)

	retNum, err = OrderNumberFromString("a")
	assert.NotNil(err)
	assert.Equal(int64(0), retNum)
	valid = OrderNumberIsValid("aaaa")
	assert.False(valid)
	valid = OrderNumberIsValid("0725-1726-4071-2412")
	assert.False(valid)

	num = OrderNumber()
	assert.True(num > 0)
	//fmt.Printf("%v\n", num)
}

func BenchmarkOrderNumber(b *testing.B) {
	for i := 0; i < 10000; i++ {
		OrderNumber()
	}
}
