package helpers

// GCD return the greatest common divisor between two ints
// Using Euclid algorithm
func GCD(a, b int) int {
	for b != 0 {
		tmp := b
		b = a % b
		a = tmp
	}
	return a
}

// Abs returns the absolute value from an int
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
