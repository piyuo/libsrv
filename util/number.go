package util

import (
	"fmt"

	"github.com/pkg/errors"
)

// ToInt convert interface{} to int
//
func ToInt(unk interface{}) (int, error) {
	switch i := unk.(type) {
	case float64:
		return int(i), nil
	case float32:
		return int(i), nil
	case int64:
		return int(i), nil
	case int32:
		return int(i), nil
	case int:
		return i, nil
	case uint64:
		return int(i), nil
	case uint32:
		return int(i), nil
	case uint:
		return int(i), nil
	default:
		return 0, errors.New(fmt.Sprintf("%v == %T\n", i, i))
	}
}

// ToFloat64 convert interface{} to int
//
func ToFloat64(unk interface{}) (float64, error) {
	switch i := unk.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint:
		return float64(i), nil
	default:
		return 0, errors.New(fmt.Sprintf("%v == %T\n", i, i))
	}
}

// ToUint32 convert interface{} to int
//
func ToUint32(unk interface{}) (uint32, error) {
	switch i := unk.(type) {
	case float64:
		return uint32(i), nil
	case float32:
		return uint32(i), nil
	case int64:
		return uint32(i), nil
	case int32:
		return uint32(i), nil
	case int:
		return uint32(i), nil
	case uint64:
		return uint32(i), nil
	case uint32:
		return i, nil
	case uint:
		return uint32(i), nil
	default:
		return 0, errors.New(fmt.Sprintf("%v == %T\n", i, i))
	}
}
