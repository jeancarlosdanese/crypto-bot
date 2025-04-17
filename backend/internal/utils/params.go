// internal/utils/params.go

package utils

func GetFloatParam(params map[string]any, key string, defaultVal float64) float64 {
	if val, ok := params[key]; ok {
		if f, ok := val.(float64); ok {
			return f
		}
	}
	return defaultVal
}

func GetIntParam(params map[string]any, key string, defaultVal int) int {
	if val, ok := params[key]; ok {
		if f, ok := val.(float64); ok { // JSON decode usa float64 para números
			return int(f)
		}
	}
	return defaultVal
}
