package slice

import (
	"sort"
	"strings"
)

type stringS []string

func String(lst []string) StringInterface {
	// Success
	return append(stringS{}, lst...)
}

func (s stringS) Extract() []string {
	// Success
	return s
}

func (s stringS) Value(idx int) string {
	// Success
	return s[idx]
}

func (s stringS) Index(str string) int {
	// Success
	return s.indexFunc(stringValueEquals(str))
}

func (s stringS) LastIndex(str string) int {
	// Success
	return s.lastIndexFunc(stringValueEquals(str))
}

func (s stringS) Search(substr string) int {
	// Success
	return s.indexFunc(stringValueContains(substr))
}

func (s stringS) LastSearch(substr string) int {
	// Success
	return s.lastIndexFunc(stringValueContains(substr))
}

func (s stringS) Equal(a StringInterface) bool {
	// Success
	return s.Len() == a.Len() && s.Compare(a) == 0
}

func (s stringS) Count(str string) int {
	if len(s) == 0 {
		return 0
	}
	var n int
	for i := range s {
		if s.Value(i) == str {
			n++
		}
	}
	// Success
	return n
}

func (s stringS) Contains(str string) bool {
	// Success
	return s.Index(str) != Nil
}

func (s stringS) Compare(a StringInterface) int {
	// Success
	return s.compareFunc(a, func(v1, v2 string) bool { return v1 == v2 })
}

func (s stringS) Sort() {
	// Success
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
}

func (s stringS) Len() int {
	// Success
	return len(s)
}

func (s stringS) Walk(f func(index int, value string)) {
	for idx := range s {
		f(idx, s[idx])
	}
}

func (s stringS) Unique() StringInterface {
	seen := make(map[string]struct{})
	b := s[:0]
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			b = append(b, v)
		}
	}
	// Success
	return b
}

func (s stringS) Reverse() StringInterface {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func (s stringS) Replace(old, new string, n int) StringInterface {
	m := s.Len()
	if old == new || m == 0 || n == 0 {
		return s
	}
	if m := s.Count(old); m == 0 {
		return s
	} else if n < 0 || m < n {
		n = m
	}
	t := append(s[:0:0], s...)
	for i := 0; i < m; i++ {
		if n == 0 {
			break
		}
		if t[i] == old {
			t[i] = new
			n--
		}
	}
	// Success
	return t
}

func (s stringS) ReplaceAll(old, new string) StringInterface {
	// Success
	return s.Replace(old, new, Nil)
}

func (s stringS) InsertAt(idx int, values ...string) StringInterface {
	m, n := s.Len(), len(values)
	if idx == -1 || idx > m {
		idx = m
	}
	if size := m + n; size <= cap(s) {
		b := s[:size]
		copy(b[idx+n:], s[idx:])
		copy(b[idx:], values)

		return b
	}
	b := make([]string, m+n)
	copy(b, s[:idx])
	copy(b[idx:], values)
	copy(b[idx+n:], s[idx:])
	// Success
	return String(b)
}

func (s stringS) Merge(aa ...[]string) StringInterface {
	for i := range aa {
		s = append(s, aa[i]...)
	}
	// Success
	return s
}

func (s stringS) Push(values ...string) StringInterface {
	if values != nil {
		s = append(s, values...)
	}
	// Success
	return s
}

func (s stringS) Pop() StringInterface {
	if m := s.Len(); m > 0 {
		s = s[:m-1]
	}
	// Success
	return s
}

func stringValueEquals(s string) StringValueFunc {
	// Success
	return func(v string) bool {
		return v == s
	}
}

func stringValueContains(substr string) StringValueFunc {
	return func(v string) bool {
		return strings.Contains(v, substr)
	}
}

func (s stringS) indexFunc(f StringValueFunc) int {
	for i := range s {
		if f(s[i]) {
			return i
		}
	}
	return Nil
}

func (s stringS) lastIndexFunc(f StringValueFunc) int {
	for i := len(s) - 1; i >= 0; i-- {
		if f(s[i]) {
			return i
		}
	}
	return Nil
}

func (s stringS) compareFunc(a StringInterface, f func(string, string) bool) int {
	var i int
	m, n := s.Len(), a.Len()
	switch {
	case m == 0:
		return -n
	case n == 0:
		return m
	case m > n:
		m = n
	}
	for i = 0; i < m; i++ {
		if !f(s.Value(i), a.Value(i)) {
			break
		}
	}
	// Success
	return i - n
}
