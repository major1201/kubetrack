package goutils

// CopyMap makes a shallow copy of a map.
func CopyMap[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}
	cp := make(map[K]V, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

// MergeMap merges multiple maps into one
func MergeMap[K comparable, V any](base map[K]V, overrides ...map[K]V) map[K]V {
	if base == nil && len(overrides) == 0 {
		return nil
	}

	if base == nil {
		base = make(map[K]V)
	}

	for _, override := range overrides {
		if override == nil {
			continue
		}
		for k, v := range override {
			base[k] = v
		}
	}

	return base
}

// MergeStrIFaceMaps merge 2 map[string]any into one
func MergeStrIFaceMaps(to, from map[string]any) map[string]any {
	out := make(map[string]any, len(to))
	for k, v := range to {
		out[k] = v
	}
	for k, v := range from {
		if v, ok := v.(map[string]any); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]any); ok {
					out[k] = MergeStrIFaceMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}
