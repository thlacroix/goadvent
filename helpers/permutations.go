package helpers

// Permute implements Heap's permutation algorithm
func Permute(values []int) [][]int {
	var permutations [][]int
	var heap func(int, []int)

	heap = func(k int, s []int) {
		if k == 1 {
			tempCopy := make([]int, len(values))
			copy(tempCopy, s)
			permutations = append(permutations, tempCopy)
		} else {
			heap(k-1, s)
			for i := 0; i < k-1; i++ {
				if k%2 == 0 {
					s[i], s[k-1] = s[k-1], s[i]
				} else {
					s[0], s[k-1] = s[k-1], s[0]
				}
				heap(k-1, s)
			}

		}
	}
	heap(len(values), values)

	return permutations
}

// Factorial computes the factorial of an int
func Factorial(n int) int {
	res := 1
	for i := 1; i <= n; i++ {
		res *= i
	}
	return res
}
