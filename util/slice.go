package util

func IsMemberOf[T comparable](ts []T, t T) bool {
	for _, i := range ts {
		if i == t {
			return true
		}
	}
	return false
}
