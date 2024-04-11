package convertor

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// ToBool convert string to boolean.
func ToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// ToString convert value to string
// for number, string, []byte, will convert to string
// for other type (slice, map, array, struct) will call json.Marshal.
// Play: https://go.dev/play/p/nF1zOOslpQq
func ToString(value any) string {
	if value == nil {
		return ""
	}

	switch val := value.(type) {
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int:
		return strconv.FormatInt(int64(val), 10)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case string:
		return val
	case []byte:
		return string(val)
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return ""
		}
		return string(b)

		// todo: maybe we should't supprt other type conversion
		// v := reflect.ValueOf(value)
		// log.Panicf("Unsupported data type: %s ", v.String())
		// return ""
	}
}

// ToJson convert value to a json string.
// Play: https://go.dev/play/p/2rLIkMmXWvR
func ToJson(value any) (string, error) {
	result, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func ToFloat(value any) (float64, error) {
	v := reflect.ValueOf(value)

	result := 0.0
	err := fmt.Errorf("ToInt: unvalid interface type %T", value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		result = float64(v.Int())
		return result, nil
	case uint, uint8, uint16, uint32, uint64:
		result = float64(v.Uint())
		return result, nil
	case float32, float64:
		result = v.Float()
		return result, nil
	case string:
		result, err = strconv.ParseFloat(v.String(), 64)
		if err != nil {
			result = 0.0
		}
		return result, err
	default:
		return result, err
	}
}

// ToInt convert value to int64 value, if input is not numerical, return 0 and error.
// Play: https://go.dev/play/p/9_h9vIt-QZ_b
func ToInt(value any) (int64, error) {
	v := reflect.ValueOf(value)

	var result int64
	err := fmt.Errorf("ToInt: invalid value type %T", value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		result = v.Int()
		return result, nil
	case uint, uint8, uint16, uint32, uint64:
		result = int64(v.Uint())
		return result, nil
	case float32, float64:
		result = int64(v.Float())
		return result, nil
	case string:
		result, err = strconv.ParseInt(v.String(), 0, 64)
		if err != nil {
			result = 0
		}
		return result, err
	default:
		return result, err
	}
}

// EncodeByte encode data to byte slice.
// Play: https://go.dev/play/p/DVmM1G5JfuP
func EncodeByte(data any) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// DecodeByte decode byte slice data to target object.
// Play: https://go.dev/play/p/zI6xsmuQRbn
func DecodeByte(data []byte, target any) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(target)
}
