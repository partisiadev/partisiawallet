package wallet

import (
	"encoding/gob"
	"errors"
)

type State int

var (
	ErrDBNotOpened      = errors.New("database not yet opened")
	ErrDBAlreadyOpened  = errors.New("database is already open")
	ErrPasswordMismatch = errors.New("password mismatch")
	ErrPasswordInvalid  = errors.New("password invalid")
	ErrPasswordNotSet   = errors.New("password not set")
)

type AppDBState int

const (
	AppDBStateIdle = AppDBState(iota)
	AppDBStateConnecting
	AppDBStateError
	AppDBStateConnected
)

type Service interface {
	Account() (Account, error)
	Accounts() ([]Account, error)
	Contacts(accountPublicKey string, offset, limit int) ([]Contact, error)
	AddUpdateAccount(account *Account) error
	AddUpdateContact(contact *Contact) (err error)
	AccountExists(publicKey string) (bool, error)
	DeleteAccounts([]Account) error
	DeleteContacts(accountPublicKey string, contacts []Contact) (int64, error)
	ContactsCount(addrPublicKey string) (int64, error)
	IsOpen() bool
	VerifyPassword(passwd string) error
}

func init() {
	gob.Register(Account{})
	gob.Register(Contact{})
}
