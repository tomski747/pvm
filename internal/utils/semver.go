package utils

import (
	"strconv"
	"strings"
)

// SemverGreater reports whether v1 is semantically greater than v2.
// Segments are compared numerically left to right; a longer version wins when
// all shared segments are equal (e.g. "3.1.1" > "3.1").
func SemverGreater(v1, v2 string) bool {
	p1 := strings.Split(v1, ".")
	p2 := strings.Split(v2, ".")
	for k := 0; k < len(p1) && k < len(p2); k++ {
		n1, _ := strconv.Atoi(p1[k])
		n2, _ := strconv.Atoi(p2[k])
		if n1 != n2 {
			return n1 > n2
		}
	}
	return len(p1) > len(p2)
}
