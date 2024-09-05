package util

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Order(a, b int) (int, int) {
	if a > b {
		return b, a
	}
	return a, b
}

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
