package util

func ToMap[K comparable, V any](s []V, keyF func(V) K) map[K]V {
	m := make(map[K]V)
	for _, v := range s {
		m[keyF(v)] = v
	}
	return m
}
