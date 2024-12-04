package functools

func Reverse[T int16 | int32 | int64 | uint | uint32 | uint64](slice []T) []T {
	reversed := make([]T, len(slice))
	for index, value := range slice {
		reversed[len(slice)-1-index] = value
	}
	return reversed
}
