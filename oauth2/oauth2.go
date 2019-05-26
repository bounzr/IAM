package oauth2

import (
	"strings"
)

// generic/util functions
func getCommonScope(scope1 string, scope2 string) (commonScope string) {
	m1 := getScopesMap(scope1)
	s2 := strings.Split(strings.TrimSpace(scope2), " ")
	commonScope = ""
	first := true
	for _, v := range s2 {
		if _, ok := m1[v]; ok {
			if first == true {
				commonScope = v
				first = false
			} else {
				commonScope = commonScope + " " + v
			}
		}
	}
	return
}

func getScopesMap(scope string) map[string]struct{} {
	m := make(map[string]struct{})
	scp := strings.Split(strings.TrimSpace(scope), " ")
	for _, v := range scp {
		m[v] = struct{}{}
	}
	return m
}
