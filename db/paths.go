package db

import (
	"encoding/gob"
)

/* Database keys to values badger mapping. each path is separated by forward slash (/)
/accounts ---> Returns all accounts
/accounts/:id ---> Returns an account with eth address or custodial id from :id params
/accounts/:id/contacts ---> Similar to above
/accounts/:id/contacts/:id
/activeAccount returns active account
*/

const (
	PathActiveAccount = `/activeAccount`
)

const (
	PathAccounts = `/accounts`
)

func init() {
	gob.Register(Account{})
	gob.Register(account{})
}
