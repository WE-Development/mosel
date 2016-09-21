package commons

func ContainsStr(arr []string, val string) bool {
	for _, entry := range arr {
		if entry == val {
			return true
		}
	}
	return false
}
