package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil { // value = nil, type = *customError, поэтому возвращает true
		println("error")
		return
	}
	println("ok")
}
