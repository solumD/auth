package sys

import (
	"github.com/pkg/errors"
	"github.com/solumD/auth/internal/sys/codes"
)

type commonError struct {
	msg  string
	code codes.Code
}

// NewCommonError ...
func NewCommonError(msg string, code codes.Code) *commonError {
	return &commonError{msg, code}
}

// Error ...
func (r *commonError) Error() string {
	return r.msg
}

// Code ...
func (r *commonError) Code() codes.Code {
	return r.code
}

// IsCommonError ...
func IsCommonError(err error) bool {
	var ce *commonError
	return errors.As(err, &ce)
}

// GetCommonError ...
func GetCommonError(err error) *commonError {
	var ce *commonError
	if !errors.As(err, &ce) {
		return nil
	}

	return ce
}
