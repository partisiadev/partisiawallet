package db

type Wallet struct {
	database *database
}
type Type string

const (
	TypeCustodial    = Type("Custodial")
	TypeNonCustodial = Type("NonCustodial")
)

var walletInstance = &Wallet{}

func Instance() *Wallet {
	if walletInstance.database == nil {
		walletInstance.database = dBInstance()
	}
	return walletInstance
}

func (d *Wallet) DBAccessor() *Accessor {
	return d.database.state
}

func (d *Wallet) AutoCreateEcdsaAccount() error {
	pvtAcc, err := autoCreateECDSAAccount()
	if err != nil {
		return err
	}
	err = d.database.saveAccount(pvtAcc)
	return nil
}
func (d *Wallet) ImportECDSAAccount(pvtKeyStr string) error {
	pvtAcc, err := importECDSAAccount(pvtKeyStr)
	if err != nil {
		return err
	}
	err = d.database.saveAccount(pvtAcc)
	return nil
}

func (d *Wallet) SetActiveAccount(account Account) error {
	return d.database.setActiveAccount(account)
}

func (d *Wallet) DeleteActiveAccount() error {
	return d.database.deleteActiveAccount()
}

func (d *Wallet) ActiveAccount() (*Account, error) {
	return d.database.activeAccount()
}

func (d *Wallet) Accounts() ([]Account, error) {
	return d.database.accounts()
}
