package utils

// Contains check if item in list
func Contains(list []interface{}, item interface{}) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
