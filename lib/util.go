package lib

func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}
