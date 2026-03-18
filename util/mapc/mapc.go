package mapc

import (
	"encoding/json"
	"math"
	"time"

	opt "github.com/zakiverse/zakiverse-api/util/optional"
)

func String(m map[string]any, key string) opt.Optional[string] {
	val, ok := m[key]
	if !ok || val == nil {
		return opt.Undefined[string]()
	}

	switch v := val.(type) {
	case string:
		return opt.Defined(v)
	default:
		return opt.Undefined[string]()
	}
}

func StringPtr(m map[string]any, key string) opt.Optional[*string] {
	val, ok := m[key]
	if !ok {
		return opt.Undefined[*string]()
	}
	if val == nil {
		return opt.Defined[*string](nil)
	}

	switch v := val.(type) {
	case string:
		return opt.Defined(opt.Ptr(v))
	default:
		return opt.Defined[*string](nil)
	}
}

func Int(m map[string]any, key string) opt.Optional[int64] {
	val, ok := m[key]
	if !ok || val == nil {
		return opt.Undefined[int64]()
	}

	switch v := val.(type) {
	case int:
		return opt.Defined(int64(v))
	case int8:
		return opt.Defined(int64(v))
	case int16:
		return opt.Defined(int64(v))
	case int32:
		return opt.Defined(int64(v))
	case int64:
		return opt.Defined(v)

	case uint:
		return opt.Defined(int64(v))
	case uint8:
		return opt.Defined(int64(v))
	case uint16:
		return opt.Defined(int64(v))
	case uint32:
		return opt.Defined(int64(v))
	case uint64:
		if v > math.MaxInt64 {
			return opt.Undefined[int64]()
		}
		return opt.Defined(int64(v))

	case float32:
		return opt.Defined(int64(v))
	case float64:
		return opt.Defined(int64(v))

	default:
		return opt.Undefined[int64]()
	}
}

func IntPtr(m map[string]any, key string) opt.Optional[*int64] {
	val, ok := m[key]
	if !ok {
		return opt.Undefined[*int64]()
	}
	if val == nil {
		return opt.Defined[*int64](nil)
	}

	v := Int(m, key)
	if opt.IsDefined(v) {
		return opt.Defined(opt.Ptr(v.V))
	}

	return opt.Defined[*int64](nil)
}

func Float(m map[string]any, key string) opt.Optional[float64] {
	val, ok := m[key]
	if !ok || val == nil {
		return opt.Undefined[float64]()
	}

	switch v := val.(type) {
	case float32:
		return opt.Defined(float64(v))
	case float64:
		return opt.Defined(v)

	case int:
		return opt.Defined(float64(v))
	case int8:
		return opt.Defined(float64(v))
	case int16:
		return opt.Defined(float64(v))
	case int32:
		return opt.Defined(float64(v))
	case int64:
		return opt.Defined(float64(v))

	case uint:
		return opt.Defined(float64(v))
	case uint8:
		return opt.Defined(float64(v))
	case uint16:
		return opt.Defined(float64(v))
	case uint32:
		return opt.Defined(float64(v))
	case uint64:
		return opt.Defined(float64(v))

	default:
		return opt.Undefined[float64]()
	}
}

func FloatPtr(m map[string]any, key string) opt.Optional[*float64] {
	val, ok := m[key]
	if !ok {
		return opt.Undefined[*float64]()
	}
	if val == nil {
		return opt.Defined[*float64](nil)
	}

	v := Float(m, key)
	if opt.IsDefined(v) {
		return opt.Defined(opt.Ptr(v.V))
	}

	return opt.Defined[*float64](nil)
}

func Boolean(m map[string]any, key string) opt.Optional[bool] {
	if val, ok := m[key]; ok {
		if cov, covok := val.(bool); covok {
			return opt.Defined(cov)
		}
	}
	return opt.Undefined[bool]()
}

func BooleanPtr(m map[string]any, key string) opt.Optional[*bool] {
	if val, ok := m[key]; ok {
		if cov, covok := val.(bool); covok {
			return opt.Defined(opt.Ptr(cov))
		}
		return opt.Defined[*bool](nil)
	}
	return opt.Undefined[*bool]()
}

func Time(m map[string]any, key string) opt.Optional[time.Time] {
	val, ok := m[key]
	if !ok || val == nil {
		return opt.Undefined[time.Time]()
	}

	switch v := val.(type) {
	case time.Time:
		return opt.Defined(v)
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return opt.Undefined[time.Time]()
		}
		return opt.Defined(t)
	default:
		return opt.Undefined[time.Time]()
	}
}

func TimePtr(m map[string]any, key string) opt.Optional[*time.Time] {
	val, ok := m[key]
	if !ok {
		return opt.Undefined[*time.Time]()
	}
	if val == nil {
		return opt.Defined[*time.Time](nil)
	}

	switch v := val.(type) {
	case time.Time:
		return opt.Defined(opt.Ptr(v))
	case *time.Time:
		return opt.Defined(v)
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return opt.Defined[*time.Time](nil)
		}
		return opt.Defined(opt.Ptr(t))
	default:
		return opt.Undefined[*time.Time]()
	}
}

func RawJson(m map[string]any, key string) opt.Optional[json.RawMessage] {
	val, ok := m[key]
	if !ok {
		return opt.Undefined[json.RawMessage]()
	}
	if val == nil {
		return opt.Defined(json.RawMessage(nil))
	}

	switch v := val.(type) {
	case json.RawMessage:
		return opt.Defined(v)
	case []byte:
		return opt.Defined(json.RawMessage(v))
	default:
		raw, err := json.Marshal(v)
		if err != nil {
			return opt.Defined(json.RawMessage(nil))
		}
		return opt.Defined(json.RawMessage(raw))
	}
}

func RawJsonPtr(m map[string]any, key string) opt.Optional[*json.RawMessage] {
	val, ok := m[key]
	if !ok {
		return opt.Undefined[*json.RawMessage]()
	}
	if val == nil {
		return opt.Defined[*json.RawMessage](nil)
	}

	raw, err := json.Marshal(val)
	if err != nil {
		return opt.Defined[*json.RawMessage](nil)
	}

	return opt.Defined(opt.Ptr[json.RawMessage](raw))
}
