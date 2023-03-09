package utils

func CopySlice(s []interface{}) []interface{} {
	result := make([]interface{}, len(s))

	for i, v := range s {
		if vm, ok := v.(map[string]interface{}); ok {
			result[i] = CopyMap(vm)
		} else if vs, ok := v.([]interface{}); ok {
			result[i] = CopySlice(vs)
		} else {
			result[i] = v
		}
	}
	return result
}

func CopyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = CopyMap(vm)
		} else if vs, ok := v.([]interface{}); ok {
			cp[k] = CopySlice(vs)
		} else {
			cp[k] = v
		}
	}
	return cp
}

func ConvertStringField(m map[string]interface{}, old string, new string) {
	if m[old] != nil {
		if m[old].(string) != "" {
			m[new] = m[old]
		}

		delete(m, old)
	}
}

func ConvertMapField(m map[string]interface{}, old string, new string) {
	if m[old] != nil {
		if len(m[old].(map[string]interface{})) > 0 {
			m[new] = m[old]
		}

		delete(m, old)
	}
}

func ConvertBoolField(m map[string]interface{}, old string, new string) {
    if m[old] != nil {
        m[new] = m[old]
        delete(m, old)
    }
}
