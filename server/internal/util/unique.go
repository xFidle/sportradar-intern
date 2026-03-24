package util

func AreUnique[T comparable](items []T) bool {
	lookup := make(map[T]struct{})
	for _, item := range items {
		if _, ok := lookup[item]; ok {
			return false
		}
		lookup[item] = struct{}{}
	}
	return true
}
