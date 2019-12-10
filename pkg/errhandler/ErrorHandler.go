package errhandler

import "fmt"

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicOnErrorWithMessage(msg string, err error) {
	if err != nil {
		fmt.Println(msg)
		panic(err)
	}
}

type ClosableWithError interface {
	Close() error
}

func CloseWithPanic(closable ClosableWithError) {
	HandleError(closable.Close())
}