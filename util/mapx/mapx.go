package mapx

import (
	"encoding/json"
	"fmt"
)

func Delete(m map[string]any, keys ...string) {
	for _, k := range keys {
		delete(m, k)
	}
}

func String(m map[string]any, key string) string {
	if val, ok := m[key]; ok {
		if cov, covok := val.(string); covok {
			return cov
		}
	}
	return ""
}

func StringExist(m map[string]any, key string) (string, bool) {
	if val, ok := m[key]; ok {
		if cov, covok := val.(string); covok {
			return cov, true
		}
		return "", true
	}
	return "", false
}

func NullableString(m map[string]any, key string) *string {
	if val, ok := m[key]; ok {
		if cov, covok := val.(string); covok {
			return &cov
		}
	}
	return nil
}

func NullableStringExist(m map[string]any, key string) (*string, bool) {
	if val, ok := m[key]; ok {
		if cov, covok := val.(string); covok {
			return &cov, true
		}
		return nil, true
	}
	return nil, false
}

func Int(m map[string]any, key string) int {
	if val, ok := m[key]; ok {
		if cov, covok := val.(float64); covok {
			return int(cov)
		}
	}
	return 0
}

func IntExist(m map[string]any, key string) (int, bool) {
	if val, ok := m[key]; ok {
		if cov, covok := val.(float64); covok {
			return int(cov), true
		}
		return 0, true
	}
	return 0, false
}

func NullableInt(m map[string]any, key string) *int {
	if val, ok := m[key]; ok {
		if cov, covok := val.(float64); covok {
			v := int(cov)
			return &v
		}
	}
	return nil
}

func NullableIntExist(m map[string]any, key string) (*int, bool) {
	if val, ok := m[key]; ok {
		if cov, covok := val.(float64); covok {
			v := int(cov)
			return &v, true
		}
		return nil, true
	}
	return nil, false
}

func Float(m map[string]any, key string) float64 {
	if val, ok := m[key]; ok {
		if cov, covok := val.(float64); covok {
			return cov
		}
	}
	return 0.0
}

func FloatExist(m map[string]any, key string) (float64, bool) {
	if val, ok := m[key]; ok {
		if cov, covok := val.(float64); covok {
			return cov, true
		}
		return 0.0, true
	}
	return 0.0, false
}

func NullableFloat(m map[string]any, key string) *float64 {
	if val, ok := m[key]; ok {
		if cov, covok := val.(float64); covok {
			return &cov
		}
	}
	return nil
}

func NullableFloatExist(m map[string]any, key string) (*float64, bool) {
	if val, ok := m[key]; ok {
		if cov, covok := val.(float64); covok {
			return &cov, true
		}
		return nil, true
	}
	return nil, false
}

func Boolean(m map[string]any, key string) bool {
	if val, ok := m[key]; ok {
		if cov, covok := val.(bool); covok {
			return cov
		}
	}
	return false
}

func BooleanExist(m map[string]any, key string) (bool, bool) {
	if val, ok := m[key]; ok {
		if cov, covok := val.(bool); covok {
			return cov, true
		}
		return false, true
	}
	return false, false
}

func NullableBoolean(m map[string]any, key string) *bool {
	if val, ok := m[key]; ok {
		if cov, covok := val.(bool); covok {
			return &cov
		}
	}
	return nil
}

func NullableBooleanExist(m map[string]any, key string) (*bool, bool) {
	if val, ok := m[key]; ok {
		if cov, covok := val.(bool); covok {
			return &cov, true
		}
		return nil, true
	}
	return nil, false
}

func Values[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

func Keys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

func ToStringMap(in map[string]any) map[string]string {
	out := make(map[string]string, len(in))

	for k, v := range in {
		switch t := v.(type) {
		case string:
			out[k] = t
		case []byte:
			out[k] = string(t)
		default:
			b, err := json.Marshal(v)
			if err != nil {
				out[k] = fmt.Sprint(v)
			} else {
				out[k] = string(b)
			}
		}
	}

	return out
}

func RawMessageToStringMap(raw *json.RawMessage) (map[string]string, error) {
	if raw == nil {
		return nil, nil
	}

	var obj map[string]any
	if err := json.Unmarshal(*raw, &obj); err != nil {
		return nil, err
	}

	return ToStringMap(obj), nil
}

func StructToMap(v any) map[string]any {
	b, _ := json.Marshal(v)
	out := map[string]any{}
	_ = json.Unmarshal(b, &out)
	return out
}
