package helper

func MapKeys(m map[int]string) (res []int) {
	for k := range m {
		res = append(res, k)
	}
	return
}

func MapKeysString(m map[string]string) (res []string) {
	for k := range m {
		res = append(res, k)
	}
	return
}
