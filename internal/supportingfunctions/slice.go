package supportingfunctions

// SliceContainsElement проверяет наличие элемента в срезе
func SliceContainsElement[T comparable](elem T, list []T) (int, bool) {
	for k, v := range list {
		if v == elem {
			return k, true
		}
	}

	return -1, false
}

// SliceContainsElementFunc проверяет наличие элемента в срезе
func SliceContainsElementFunc[L any](list []L, f func(n int) bool) (int, bool) {
	for k := range list {
		if f(k) {
			return k, true
		}
	}

	return -1, false
}
