package domain

import "fmt"

type DomainError struct {
	Code    ErrorCode
	Message string
}

func (e DomainError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return string(e.Code)
}

func NewDomainError(code ErrorCode, format string, args ...any) error {
	return DomainError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}



