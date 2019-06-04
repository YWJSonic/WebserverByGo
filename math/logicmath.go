package math

// IsInclude ...
func IsInclude(target int, src []int) bool {
	for _, value := range src {
		if value == target {
			return true
		}
	}
	return false
}
