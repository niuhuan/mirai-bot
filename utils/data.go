package utils

func ContainsInt64(items []int64, item int64) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func ContainsInt(items []int, item int) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}


func SumInts(items []int) (sum int) {
	for _, item := range items {
		sum += item
	}
	return
}
