package support

func Tap(value interface{}, callback func(...interface{})) interface{} {
	callback(value)

	return value
}
