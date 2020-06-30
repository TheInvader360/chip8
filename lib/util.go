package lib

func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Clamp(n int, min int, max int) int {
	return Max(min, Min(max, n))
}
