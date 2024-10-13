package util

// Allocates a new object and assigns the value passed in the parameter.
func New[T any](val T) *T {
	res := new(T)
	*res = val

	return res
}
