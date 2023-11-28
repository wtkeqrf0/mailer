package ptr

func Get[T any](t T) *T {
	return &t
}
