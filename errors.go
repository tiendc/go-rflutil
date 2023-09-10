package rflutil

import "errors"

var (
	ErrTypeInvalid        = errors.New("ErrTypeInvalid")
	ErrTypeUnmatched      = errors.New("ErrTypeUnmatched")
	ErrNotFound           = errors.New("ErrNotFound")
	ErrValueUnaddressable = errors.New("ErrValueUnaddressable")
	ErrValueUnsettable    = errors.New("ErrValueUnsettable")
	ErrIndexOutOfRange    = errors.New("ErrIndexOutOfRange")
)
