package main

func panicT[T any](v T, err error) T {
	panicIf(err)
	return v
}
func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
