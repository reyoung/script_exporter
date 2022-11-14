package main

func mapKeys[k comparable, v any](m map[k]v) []k {
	if len(m) == 0 {
		return nil
	}
	ret := make([]k, 0, len(m))
	for key := range m {
		ret = append(ret, key)
	}
	return ret
}
