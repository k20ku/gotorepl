package gotorepl

import "strings"

func hasPrefix(target string, rs ...rune) bool {
	res := false
	for _, r := range rs {
		res = res || strings.HasPrefix(strings.TrimSpace(target), string(r))
	}
	return res
}
