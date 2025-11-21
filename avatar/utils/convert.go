package utils

import "strconv"

func ConvertToInt64(val interface{}) int64 {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
		return 0
	case float64:
		return int64(v)
	default:
		return 0
	}
}
