package math

// Abs for int64
func Abs(value int64) int64 {
	if value < 0 {
		return value * -1
	}
	return value

}
