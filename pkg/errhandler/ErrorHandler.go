package errhandler

import "fmt"

func PanicOnError(err error) {
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

func CloseWithPanicOnError(closable ClosableWithError) {
	PanicOnError(closable.Close())
}