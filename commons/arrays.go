package commons

func ContainsStr(arr []string, val string) bool {
	for _, entry := range arr {
		if entry == val {
			return true
		}
	}
	return false
}

func ExcludeStr(arr []string, ex []string) []string {
	res := make([]string, 0)

	elems: for _, elem := range arr {
		for _, ex := range ex {
			if elem == ex {
				continue elems
			}
		}
		res = append(res, elem)
	}

	return res
}
