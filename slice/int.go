package slice

import (
	"sort"
)

type intS []int

func Int(lst []int) IntInterface {
	// Success
	return append(intS{}, lst...)
}

func (s intS) Extract() []int {
	// Success
	return s
}

func (s intS) Value(idx int) int {
	// Success
	return s[idx]
}

func (s intS) Index(i int) int {
	// Success
	return s.indexFunc(intValueEquals(i))
}

func (s intS) LastIndex(i int) int {
	// Success
	return s.lastIndexFunc(intValueEquals(i))
}

func (s intS) Equal(a IntInterface) bool {
	// Success
	return s.Len() == a.Len() && s.Compare(a) == 0
}

func (s intS) Count(i int) int {
	if len(s) == 0 {
		return 0
	}
	var n int
	for i := range s {
		if s.Value(i) == i {
			n++
		}
	}
	// Success
	return n
}

func (s intS) Contains(i int) bool {
	// Success
	return s.Index(i) != Nil
}

func (s intS) Compare(a IntInterface) int {
	// Success
	return s.compareFunc(a, func(v1, v2 int) bool { return v1 == v2 })
}

func (s intS) Sort() {
	// Success
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
}

func (s intS) Len() int {
	// Success
	return len(s)
}

func (s intS) Walk(f func(index int, value int)) {
	for idx := range s {
		f(idx, s[idx])
	}
}

func (s intS) Unique() IntInterface {
	seen := make(map[int]struct{})
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

func (s intS) Reverse() IntInterface {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func (s intS) Replace(old, new int, n int) IntInterface {
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

func (s intS) ReplaceAll(old, new int) IntInterface {
	// Success
	return s.Replace(old, new, Nil)
}

func (s intS) InsertAt(idx int, values ...int) IntInterface {
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
	b := make([]int, m+n)
	copy(b, s[:idx])
	copy(b[idx:], values)
	copy(b[idx+n:], s[idx:])
	// Success
	return Int(b)
}

func (s intS) Merge(aa ...[]int) IntInterface {
	for i := range aa {
		s = append(s, aa[i]...)
	}
	// Success
	return s
}

func (s intS) Push(values ...int) IntInterface {
	if values != nil {
		s = append(s, values...)
	}
	// Success
	return s
}

func (s intS) Pop() IntInterface {
	if m := s.Len(); m > 0 {
		s = s[:m-1]
	}
	// Success
	return s
}

func intValueEquals(i int) IntValueFunc {
	// Success
	return func(v int) bool {
		return v == i
	}
}

func (s intS) indexFunc(f IntValueFunc) int {
	for i := range s {
		if f(s[i]) {
			return i
		}
	}
	return Nil
}

func (s intS) lastIndexFunc(f IntValueFunc) int {
	for i := len(s) - 1; i >= 0; i-- {
		if f(s[i]) {
			return i
		}
	}
	return Nil
}

func (s intS) compareFunc(a IntInterface, f func(int, int) bool) int {
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
