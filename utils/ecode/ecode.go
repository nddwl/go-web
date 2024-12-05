package ecode

import (
	"errors"
	"fmt"
	"strconv"
)

var _codes = map[int]struct{}{}

func New(code int, msg string) Code {
	_, ok := _codes[code]
	if ok {
		panic("存在相同code")
	}
	_codes[code] = struct{}{}
	return Code{
		code:    code,
		message: msg,
	}
}

type Codes interface {
	Error() string
	Code() int
	Message() string
}

var _ Codes = Code{}

type Code struct {
	code    int
	message string
}

func (e Code) Error() string {
	return strconv.Itoa(e.code) + e.message
}

func (e Code) Code() int {
	return e.code
}

func (e Code) Message() string {
	return e.message
}

func Cause(err error) Codes {
	if err == nil {
		return OK
	}
	ecode := Code{}
	ok := errors.As(err, &ecode)
	if ok {
		return ecode
	}
	fmt.Println("异常错误:", err)
	return ServerErr
}
