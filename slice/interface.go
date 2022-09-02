package slice

type (
	StringValueFunc func(v string) bool
	StringInterface interface {
		Extract() []string
		Value(idx int) string
		Index(str string) int
		LastIndex(str string) int
		Search(substr string) int
		LastSearch(substr string) int
		Equal(a StringInterface) bool
		Count(str string) int
		Contains(str string) bool
		Compare(a StringInterface) int
		Sort()
		Len() int
		Walk(f func(index int, value string))
		Unique() StringInterface
		Reverse() StringInterface
		Replace(old, new string, n int) StringInterface
		ReplaceAll(old, new string) StringInterface
		InsertAt(idx int, values ...string) StringInterface
		Merge(aa ...[]string) StringInterface
		Push(values ...string) StringInterface
		Pop() StringInterface
	}

	IntValueFunc func(v int) bool
	IntInterface interface {
		Extract() []int
		Value(idx int) int
		Index(i int) int
		LastIndex(i int) int
		Equal(a IntInterface) bool
		Count(i int) int
		Contains(i int) bool
		Compare(a IntInterface) int
		Sort()
		Len() int
		Walk(f func(index int, value int))
		Unique() IntInterface
		Reverse() IntInterface
		Replace(old, new int, n int) IntInterface
		ReplaceAll(old, new int) IntInterface
		InsertAt(idx int, values ...int) IntInterface
		Merge(aa ...[]int) IntInterface
		Push(values ...int) IntInterface
		Pop() IntInterface
	}

	Int64ValueFunc func(v int64) bool
	Int64Interface interface {
		Extract() []int64
		Value(idx int) int64
		Index(i int64) int
		LastIndex(i int64) int
		Equal(a Int64Interface) bool
		Count(i int64) int
		Contains(i int64) bool
		Compare(a Int64Interface) int
		Sort()
		Len() int
		Walk(f func(index int, value int64))
		Unique() Int64Interface
		Reverse() Int64Interface
		Replace(old, new int64, n int) Int64Interface
		ReplaceAll(old, new int64) Int64Interface
		InsertAt(idx int, values ...int64) Int64Interface
		Merge(aa ...[]int64) Int64Interface
		Push(values ...int64) Int64Interface
		Pop() Int64Interface
	}
)

var (
	Nil = -1
)
