package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"gioui.org/app"
	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/partisiadev/partisiawallet/app/assets/evm"
	"github.com/partisiadev/partisiawallet/config"
	"github.com/partisiadev/partisiawallet/log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var GlobalWallet = NewWallet()

// const defaultChainMediator = "https://mainnet.infura.io"
// const defaultChainMediator = "https://ethereum.publicnode.com"
const defaultChainMediator = "https://rpc.ntity.io"

type RpcClient struct {
	ApiKey string
	RPC    evm.RPC
}

type Wallet struct {
	//connections    map[string][]RpcClient
	password       string
	stateInternal  appDBStateInternal
	state          AppDBState
	dB             *badger.DB
	err            error
	FavoriteChains map[string]struct{}
	FavoriteRPCs   map[string]struct{}
}
type appDBStateInternal struct {
	state AppDBState
	dB    *badger.DB
	err   error
}

func NewWallet() *Wallet {
	wa := &Wallet{}
	wa.FavoriteChains = make(map[string]struct{})
	wa.FavoriteRPCs = make(map[string]struct{})
	return wa
}

func (w *Wallet) CreateAccount(pvtKeyStr string) (err error) {
	if strings.TrimSpace(pvtKeyStr) == "" {
		err = errors.New("private key is empty")
		return
	}
	if !w.IsOpen() {
		return ErrDBNotOpened
	}
	privateKey, err := crypto.HexToECDSA(pvtKeyStr)
	if err != nil {
		return err
	}
	publicKey := privateKey.PublicKey
	publicKeyStr := hex.EncodeToString(crypto.FromECDSAPub(&publicKey))
	ethAddress := hex.EncodeToString(crypto.PubkeyToAddress(publicKey).Bytes())

	account := Account{
		PrivateKey: pvtKeyStr,
		PublicKey:  publicKeyStr,
		EthAddress: ethAddress,
	}
	err = w.AddUpdateAccount(&account)
	return err
}

func (w *Wallet) AutoCreateAccount() (err error) {
	defer func() {
		if err != nil {
			log.Logger().Errorln(err)
		}
	}()
	ecdsaPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}

	publicKey := ecdsaPrivateKey.PublicKey
	pvtKeyStr := hex.EncodeToString(ecdsaPrivateKey.D.Bytes())
	publicKeyStr := hex.EncodeToString(crypto.FromECDSAPub(&publicKey))
	ethAddress := hex.EncodeToString(crypto.PubkeyToAddress(publicKey).Bytes())
	account := Account{
		PrivateKey: pvtKeyStr,
		PublicKey:  publicKeyStr,
		EthAddress: ethAddress,
	}
	err = w.AddUpdateAccount(&account)
	return err
}
func (w *Wallet) Open(options badger.Options) error {
	_ = w.Close()
	state := w.getState()
	state.state = AppDBStateConnecting
	defer func() {
		w.setState(state)
	}()
	dB, err := badger.Open(options)
	if err != nil {
		state.err = err
		state.dB = nil
		state.state = AppDBStateError
		return err
	}
	state.err = nil
	state.dB = dB
	state.state = AppDBStateConnected
	return nil
}

func (w *Wallet) Close() error {
	state := w.getState()
	if state.dB != nil {
		_ = w.Close()
	}
	state.err = nil
	state.dB = nil
	state.state = AppDBStateIdle
	w.setState(state)
	return nil
}

func (w *Wallet) getState() appDBStateInternal {
	return w.stateInternal
}
func (w *Wallet) setState(state appDBStateInternal) {
	w.stateInternal = state
}

func (w *Wallet) getErrorState() (err error) {
	state := w.getState()
	if state.dB == nil {
		return ErrDBNotOpened
	}
	if state.err != nil {
		return state.err
	}
	return nil
}

// prefixScan prefixPos is the position of separator to derive prefix key for scanning
func (w *Wallet) prefixScan(prefixOrFullKey string, keySeparator string, prefixPos int) (fullKeys []string, err error) {
	err = w.getErrorState()
	if err != nil {
		return fullKeys, err
	}
	db := w.getState().dB
	prefixKeyArr := strings.Split(prefixOrFullKey, keySeparator)
	if len(prefixKeyArr) < prefixPos+1 {
		return fullKeys, ErrInvalidKey
	}
	prefixKey := strings.Join(prefixKeyArr[0:prefixPos+1], keySeparator)
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek([]byte(prefixKey)); it.ValidForPrefix([]byte(prefixKey)); it.Next() {
			item := it.Item()
			k := string(item.KeyCopy(nil))
			fullKeys = append(fullKeys, k)
		}
		return nil
	})
	return fullKeys, err
}

func (w *Wallet) prefixScanSorted(prefixOrFullKey, sep string, prefixSepPos, sortSepPos int, shouldReverse bool) (derivedKeys []string, err error) {
	derivedKeys, err = w.prefixScan(prefixOrFullKey, sep, prefixSepPos)
	if err != nil {
		return
	}
	if len(derivedKeys) > 0 {
		arr := strings.Split(derivedKeys[0], KeySeparator)
		if sortSepPos >= len(arr) {
			return derivedKeys, errors.New("invalid separator pos")
		}
	}
	sort.Slice(derivedKeys, func(i, j int) bool {
		keyOneArr := strings.Split(derivedKeys[i], KeySeparator)
		keyTwoArr := strings.Split(derivedKeys[j], KeySeparator)
		keyOne := keyOneArr[sortSepPos]
		keyTwo := keyTwoArr[sortSepPos]
		if shouldReverse {
			return keyOne > keyTwo
		}
		return keyOne < keyTwo
	})
	return derivedKeys, err
}

// ViewRecord
//
//	ptrStruct should be a pointer to a struct registered with gob
func (w *Wallet) viewRecord(key []byte, ptrStruct interface{}) (err error) {
	err = w.getErrorState()
	if err != nil {
		return err
	}
	dB := w.getState().dB
	err = dB.View(func(txn *badger.Txn) (err error) {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) (err error) {
			err = DecodeToStruct(ptrStruct, val)
			return err
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (w *Wallet) IsOpen() bool {
	state := w.getState()
	if state.dB == nil {
		return false
	}
	return !state.dB.IsClosed()
}

func (w *Wallet) OpenFromPassword(passwd string) error {
	origPasswd := passwd
	if w.IsOpen() {
		return ErrDBAlreadyOpened
	}
	if len(passwd) == 0 {
		return ErrPasswdCannotBeEmpty
	}
	if len([]byte(passwd)) > MaxNumOfPasswdChars {
		return fmt.Errorf("password should be less than %w characters", MaxNumOfPasswdChars)
	}
	padDiff := MaxNumOfPasswdChars - len([]byte(passwd))
	var leftPad, rightPad string
	for i := 0; i < padDiff; i++ {
		if i%2 == 0 {
			leftPad += passwdPadCharacter
		} else {
			rightPad += passwdPadCharacter
		}
	}
	passwd = leftPad + passwd + rightPad
	dirPath, err := app.DataDir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(dirPath, config.AppConfigPathDirName, config.AppWalletDBDirName)
	options := badger.DefaultOptions(dbPath)
	options.EncryptionKey = []byte(passwd)
	options.IndexCacheSize = 100
	err = w.Open(options)
	if err != nil {
		st := w.getState()
		st.err = err
		w.setState(st)
		return err
	}
	w.password = origPasswd
	return nil
}
func (w *Wallet) DatabaseExists() bool {
	dirPath, err := app.DataDir()
	if err != nil {
		log.Logger().Fatal(err)
	}
	dbPath := filepath.Join(dirPath, config.AppConfigPathDirName, config.AppWalletDBDirName)
	_, err = os.Stat(dbPath)
	return err == nil
}

// VerifyPassword
// returns nil if password is correct else may
// return ErrPasswordMismatch or ErrPasswordInvalid or ErrPasswordNotSet
func (w *Wallet) VerifyPassword(passwd string) error {
	if strings.TrimSpace(passwd) == "" {
		return ErrPasswordInvalid
	}
	if w.password == "" {
		return ErrPasswordNotSet
	}
	if w.password != passwd {
		return ErrPasswordMismatch
	}
	return nil
}

// AddUpdateAccount saves as primary account
func (w *Wallet) AddUpdateAccount(acc *Account) (err error) {
	if acc == nil || len(acc.PublicKey) == 0 {
		return ErrInvalidAccount
	}
	//prevAccount, _ := w.Account()
	//var currentAccountChanged bool
	err = w.getErrorState()
	if err != nil {
		return err
	}
	dB := w.getState().dB
	defer func() {
		if r := recover(); r != nil {
			log.Logger().Errorln(r)
		}
		//if currentAccountChanged {
		//	event := pubsub.Event{Data: pubsub.CurrentAccountChangedEventData{
		//		PrevAccountPublicKey:    prevAccount.PublicKey,
		//		CurrentAccountPublicKey: acc.PublicKey,
		//	}, Topic: pubsub.CurrentAccountChangedEventTopic}
		//	w.EventBroker.Fire(event)
		//}
	}()
	acc.UpdatedAt = time.Now()
	if acc.CreatedAt.IsZero() {
		acc.CreatedAt = time.Now()
	}
	fullKey, err := acc.GetDBFullKey()
	if err != nil {
		return
	}
	var accountKeys []string
	txn := dB.NewTransaction(true)
	defer txn.Discard()
	opts := badger.DefaultIteratorOptions
	it := txn.NewIterator(opts)
	for it.Seek([]byte(KeyPrefixAccounts)); it.ValidForPrefix([]byte(KeyPrefixAccounts)); it.Next() {
		item := it.Item()
		k := string(item.KeyCopy(nil))
		accountKeys = append(accountKeys, k)
	}
	it.Close()
	for _, key := range accountKeys {
		if strings.Contains(key, acc.PublicKey) {
			err = txn.Delete([]byte(key))
		}
	}
	val := EncodeToBytes(acc)
	err = txn.Set([]byte(fullKey), val)
	if err != nil {
		return err
	}
	err = txn.Commit()
	//if err == nil && prevAccount.PublicKey != acc.PublicKey {
	//	currentAccountChanged = true
	//}
	return err
}

func (w *Wallet) Account() (acc Account, err error) {
	err = w.getErrorState()
	if err != nil {
		return acc, err
	}
	accs, err := w.Accounts()
	if err != nil {
		return acc, err
	}
	if len(accs) == 0 {
		return acc, ErrAccountDoesNotExist
	}
	return accs[0], err
}

func (w *Wallet) Accounts() (accounts []Account, err error) {
	err = w.getErrorState()
	if err != nil {
		return accounts, err
	}
	prefix := KeyPrefixAccounts
	allKeys, err := w.prefixScanSorted(prefix, KeySeparator, 0, 1, true)
	if err != nil {
		return accounts, err
	}
	for _, k := range allKeys {
		var account Account
		err = w.viewRecord([]byte(k), &account)
		if err != nil {
			return accounts, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// DeleteAccounts cascade deletes all Contacts and Messages belong to those Accounts
func (w *Wallet) DeleteAccounts(accounts []Account) (err error) {
	if len(accounts) == 0 {
		return nil
	}
	defer func() {
		//w.EventBroker.Fire(pubsub.Event{
		//	Data:  pubsub.AccountsChangedEventData{},
		//	Topic: pubsub.AccountsChangedEventTopic,
		//})
	}()

	err = w.getErrorState()
	if err != nil {
		return err
	}
	dB := w.getState().dB
	for _, eachAccount := range accounts {
		err = dB.Update(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			it := txn.NewIterator(opts)
			defer it.Close()
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				k := item.Key()
				if strings.Contains(string(k), eachAccount.PublicKey) {
					err = txn.Delete(k)
				}
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return err
}

func (w *Wallet) AccountExists(publicKey string) (exists bool, err error) {
	err = w.getErrorState()
	if err != nil {
		return exists, err
	}
	dB := w.getState().dB
	if len(publicKey) == 0 {
		return exists, ErrInvalidKey
	}
	acc := Account{PublicKey: publicKey}
	accountKey := acc.GetAccountDBPrefixKey()
	err = dB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		it.Seek([]byte(accountKey))
		exists = it.ValidForPrefix([]byte(accountKey))
		it.Close()
		return nil
	})
	return exists, err
}

func (w *Wallet) Contacts(accountPublicKey string, offset, limit int) (contacts []Contact, err error) {
	contact := Contact{AccountPublicKey: accountPublicKey}
	keyPrefix := contact.GetDBPrefixKey()
	allKeys, err := w.prefixScanSorted(keyPrefix, KeySeparator, 1, 3, true)
	if err != nil {
		return contacts, err
	}
	if len(allKeys) <= offset {
		return contacts, errors.New("invalid offset")
	}
	availableLimit := len(allKeys[offset:])
	if limit > availableLimit {
		limit = availableLimit
	}
	allKeys = allKeys[offset : offset+limit]
	for _, k := range allKeys {
		var contact Contact
		err = w.viewRecord([]byte(k), &contact)
		if err == nil {
			contacts = append(contacts, contact)
		}
	}
	return contacts, err
}

func (w *Wallet) AddUpdateContact(c *Contact) (err error) {
	err = w.getErrorState()
	if err != nil {
		return err
	}
	dB := w.getState().dB
	fullKey, err := c.GetDBFullKey()
	if err != nil {
		return err
	}
	duplicateKeys, err := w.prefixScan(fullKey, KeySeparator, 2)
	if err != nil {
		return err
	}
	txn := dB.NewTransaction(true)
	defer txn.Discard()
	if len(duplicateKeys) > 0 {
		for _, key := range duplicateKeys {
			err = txn.Delete([]byte(key))
			if err != nil {
				return err
			}
		}
	}
	c.UpdatedAt = time.Now()
	fullKey, err = c.GetDBFullKey()
	if err != nil {
		return err
	}
	val := EncodeToBytes(c)
	err = txn.Set([]byte(fullKey), val)
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	//eventData := pubsub.SaveContactEventData{Contact: *c}
	//event := pubsub.Event{
	//	Data:  eventData,
	//	Topic: pubsub.SaveContactTopic,
	//}
	//w.EventBroker.Fire(event)
	return err
}

func (w *Wallet) ContactsCount(accountPublicKey string) (count int64, err error) {
	err = w.getErrorState()
	if err != nil {
		return count, err
	}
	dB := w.getState().dB
	defer func() {
		if r := recover(); r != nil {
			log.Logger().Errorln(r)
		}
	}()
	c := Contact{AccountPublicKey: accountPublicKey}
	prefixKey := c.GetDBPrefixKey()
	// This should never happen
	if err != nil {
		return
	}
	err = dB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := string(item.Key())
			if strings.HasPrefix(k, prefixKey) {
				count++
			}
		}
		return nil
	})
	return
}

func (w *Wallet) DeleteContacts(accountPublicKey string, contacts []Contact) (count int64, err error) {
	err = w.getErrorState()
	if err != nil {
		return count, err
	}
	dB := w.getState().dB
	if len(contacts) == 0 {
		return count, errors.New("contacts is empty")
	}
	if len(accountPublicKey) == 0 {
		return count, errors.New("account public key is empty")
	}
	defer func() {
		if err != nil {
			log.Logger().Errorln(err)
		}
		if count > 0 {
			//w.EventBroker.Fire(pubsub.Event{
			//	Data: pubsub.ContactsChangeEventData{
			//		AccountPublicKey: accountPublicKey,
			//	},
			//	Topic: pubsub.ContactsChangedEventTopic,
			//})
		}
	}()
	// This will delete all the contacts and all messages that belongs to deleted contact
	for _, eachContact := range contacts {
		keyComponent := fmt.Sprintf("%s%s%s", accountPublicKey, KeySeparator, eachContact.PublicKey)
		err = dB.Update(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			it := txn.NewIterator(opts)
			defer it.Close()
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				k := item.Key()
				if strings.Contains(string(k), keyComponent) {
					err = txn.Delete(k)
					if err == nil {
						count++
					}
				}
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return count, err
		}
	}
	return count, nil
}
