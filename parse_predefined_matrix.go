package main

import "strings"

func parsePredefinedMatrix(mat string) map[string][]string {
	result := make(map[string][]string)

	if len(mat) == 0 {
		return result
	}

	for _, item := range strings.Split(mat, ":") {
		splits := strings.SplitN(item, "=", 2)
		key := splits[0]
		vals := splits[1]
		result[key] = strings.Split(vals, ",")
	}
	return result
}
