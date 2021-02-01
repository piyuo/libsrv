package identifier

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// OrderNumber is easy read but not guarantee to be unique id,collision chance is 1/10,000,000,000. it format like credit card number, e.g. 0623847666125355 , first 4 digit is current date ,5-15 is random number,last digit is check sum
//
//	id := OrderNumber() //16249128003811148
//
func OrderNumber() int64 {
	var sum int64
	var num int64
	var n int64
	t := time.Now()
	num += (int64(t.Month())*100 + int64(t.Day())) * 1000000000000
	base := int64(10)
	for i := 0; i < 11; i++ {
		n = rand.Int63n(9)
		num += n * base
		base = base * 10
		if i%2 == 0 {
			n *= 2
			// If the result of this doubling operation is greater than 9.
			if n > 9 {
				// The same final result can be found by subtracting 9 from that result.
				n -= 9
			}
		}
		sum += n
	}
	luhn := sum % 10
	if luhn != 0 {
		luhn = 10 - luhn
	}
	num += int64(luhn)
	return num
}

// OrderNumberToString convert order number to easy ready string like 0624-9128-0038-11148
//
//	id := OrderNumberToString() //0624-9128-0038-11148
//
func OrderNumberToString(num int64) string {
	str := strconv.FormatInt(num, 10)
	if len(str) == 15 {
		str = "0" + str
	}
	a := []rune(str)
	return string(a[0:4]) + "-" + string(a[4:8]) + "-" + string(a[8:12]) + "-" + string(a[12:16])
}

// OrderNumberFromString convert order number string back to number like 6249128003811148
//
//	num := OrderNumberFromString("0624-9128-0038-11148") //6249128003811148
//
func OrderNumberFromString(str string) (int64, error) {
	str = strings.ReplaceAll(str, "-", "")
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse number: "+str)
	}
	return n, nil
}

const (
	asciiZero = 48
	asciiTen  = 57
)

// OrderNumberIsValid return true if order number is valid
//
//	valid := CheckNumberIsValid("0624-9128-0038-11148") //true
//
func OrderNumberIsValid(str string) bool {
	if len(str) != 19 || !strings.Contains(str, "-") {
		return false
	}
	str = strings.ReplaceAll(str, "-", "")
	var sum int64
	for i, d := range str[4:15] {
		if d < asciiZero || d > asciiTen {
			return false
		}
		d = d - asciiZero
		// Double the value of every second digit.
		if i%2 == 0 {
			d *= 2
			// If the result of this doubling operation is greater than 9.
			if d > 9 {
				// The same final result can be found by subtracting 9 from that result.
				d -= 9
			}
		}
		// Take the sum of all the digits.
		sum += int64(d)
	}
	luhn := sum % 10
	if luhn != 0 {
		luhn = 10 - luhn
	}
	checkSum := int(str[15]) - asciiZero
	if int64(checkSum) == luhn {
		return true
	}
	return false

}
