package db

import "errors"

var (
	ErrDBIsClosed          = errors.New("database is closed")
	ErrDBIsAlreadyOpen     = errors.New("database is already open")
	ErrPasswordMismatch    = errors.New("password mismatch")
	ErrPasswordInvalid     = errors.New("password invalid")
	ErrPasswordNotSet      = errors.New("password not set")
	ErrPasswdCannotBeEmpty = errors.New("password cannot be empty")
	ErrInvalidPathID       = errors.New("path id is invalid")
	ErrInvalidPath         = errors.New("invalid path")
	ErrInvalidAccount      = errors.New("invalid account")
)
