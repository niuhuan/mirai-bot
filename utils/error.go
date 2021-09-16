package utils

func PanicNotNil(err interface{}) {
	if err != nil {
		panic(err)
	}
}
