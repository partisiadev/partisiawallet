package db

type Account struct {
	account account
}

type AccountType int

const (
	AccountTypeInvalid = AccountType(iota)
	AccountTypeNonCustodial
	AccountTypeCustodial
)

type account struct {
	Mnemonics  string
	PrivateKey string
	PublicKey  string
	// For noncustodial account eth address is used as key in DB
	EthAddress string
	// For custodial account custodial id is used as key in DB
	CustodialID   string
	CustodialData interface{}
}

func (a *Account) PathID() (pathID string) {
	pathID = a.account.EthAddress
	if pathID == "" {
		pathID = a.account.CustodialID
	}
	return pathID
}
func (a *Account) Type() AccountType {
	ethAdd := a.account.EthAddress
	custID := a.account.CustodialID
	if ethAdd == "" && custID == "" {
		return AccountTypeInvalid
	}
	if ethAdd != "" {
		return AccountTypeNonCustodial
	}
	if custID != "" {
		return AccountTypeCustodial
	}
	return AccountTypeInvalid
}
