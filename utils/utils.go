package utils

func Panic(err error, msg string) {
	if err != nil {
		print(msg)
		panic(err)
	}
}
