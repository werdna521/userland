package slice

// AnyStr is similar to JS' Array.any, for string slices
func AnyStr(slice []string, f func(string) bool) bool {
	for _, v := range slice {
		if f(v) {
			return true
		}
	}
	return false
}
