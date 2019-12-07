package customError

import (
	"errors"
	"strings"
)

var ErrorTimeOut=errors.New("TimeOut")
var ErrorBufferExhausted=errors.New("Buffer Exhausted")
var ErrorInvalidFile=errors.New("Invalid File")
var ErrorZeroUserFound=errors.New("Invalid users Found")

func IsErrorNginx429(err error) bool {
	buffer := err.Error()
	if strings.Contains(buffer, "429 Too Many Requests") {
		return true
	}
	return false
}

func IsInvalidCard(err error) bool {
	buffer := err.Error()
	if strings.Contains(buffer, "invalid_card") {
		return true
	}
	return false
}


func NoSuchHost(err error) bool {
	buffer := err.Error()
	if strings.Contains(buffer, "no such host") {
		return true
	}
	return false
}

