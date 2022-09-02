package slice

import "sort"

type int64S []int64

func Int64(lst []int64) Int64Interface {
	// Success
	return append(int64S{}, lst...)
}

func (s int64S) Extract() []int64 {
	// Success
	return s
}

func (s int64S) Value(idx int) int64 {
	// Success
	return s[idx]
}

func (s int64S) Index(i int64) int {
	// Success
	return s.indexFunc(int64ValueEquals(i))
}

func (s int64S) LastIndex(i int64) int {
	// Success
	return s.lastIndexFunc(int64ValueEquals(i))
}

func (s int64S) Equal(a Int64Interface) bool {
	// Success
	return s.Len() == a.Len() && s.Compare(a) == 0
}

func (s int64S) Count(i int64) int {
	if len(s) == 0 {
		return 0
	}
	var n int
	for item := range s {
		if s.Value(item) == i {
			n++
		}
	}
	// Success
	return n
}

func (s int64S) Contains(i int64) bool {
	// Success
	return s.Index(i) != Nil
}

func (s int64S) Compare(a Int64Interface) int {
	// Success
	return s.compareFunc(a, func(v1, v2 int64) bool { return v1 == v2 })
}

func (s int64S) Sort() {
	// Success
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
}

func (s int64S) Len() int {
	// Success
	return len(s)
}

func (s int64S) Walk(f func(index int, value int64)) {
	for idx := range s {
		f(idx, s[idx])
	}
}

func (s int64S) Unique() Int64Interface {
	seen := make(map[int64]struct{})
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

func (s int64S) Reverse() Int64Interface {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	// Success
	return s
}

func (s int64S) Replace(old, new int64, n int) Int64Interface {
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

func (s int64S) ReplaceAll(old, new int64) Int64Interface {
	// Success
	return s.Replace(old, new, Nil)
}

func (s int64S) InsertAt(idx int, values ...int64) Int64Interface {
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
	b := make([]int64, m+n)
	copy(b, s[:idx])
	copy(b[idx:], values)
	copy(b[idx+n:], s[idx:])
	// Success
	return Int64(b)
}

func (s int64S) Merge(aa ...[]int64) Int64Interface {
	for i := range aa {
		s = append(s, aa[i]...)
	}
	// Success
	return s
}

func (s int64S) Push(values ...int64) Int64Interface {
	if values != nil {
		s = append(s, values...)
	}
	// Success
	return s
}

func (s int64S) Pop() Int64Interface {
	if m := s.Len(); m > 0 {
		s = s[:m-1]
	}
	// Success
	return s
}

func int64ValueEquals(i int64) Int64ValueFunc {
	// Success
	return func(v int64) bool {
		return v == i
	}
}

func (s int64S) indexFunc(f Int64ValueFunc) int {
	for i := range s {
		if f(s[i]) {
			return i
		}
	}
	return Nil
}

func (s int64S) lastIndexFunc(f Int64ValueFunc) int {
	for i := len(s) - 1; i >= 0; i-- {
		if f(s[i]) {
			return i
		}
	}
	return Nil
}

func (s int64S) compareFunc(a Int64Interface, f func(int64, int64) bool) int {
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
