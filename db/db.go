package db

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v4"
)

var (
	ErrInvalidCollectionPath = errors.New("collection path is not valid")
	ErrInvalidDocumentPath   = errors.New("document path is not valid")
)

type database struct {
	state *Accessor
}

var databaseInstance = &database{}

func dBInstance() *database {
	if databaseInstance.state == nil {
		databaseInstance.state = accessorInstance()
	}
	return databaseInstance
}

func (d *database) accounts() ([]Account, error) {
	allKeys, err := d.loadKeysForAccounts()
	if err != nil {
		return nil, err
	}
	accs, err := d.loadAccountsFromKeys(allKeys)
	return accs, err
}

func (d *database) loadKeysForAccounts() ([]string, error) {
	db, err := d.state.getOpenedDB()
	if err != nil {
		return nil, err
	}
	allKeys := make([]string, 0)
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek([]byte(PathAccounts)); it.ValidForPrefix([]byte(PathAccounts)); it.Next() {
			item := it.Item()
			kbs := item.KeyCopy(nil)
			k := string(kbs)
			if bytes.HasSuffix(kbs, []byte("/")) {
				break
			}
			allKeys = append(allKeys, k)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return allKeys, err
}

func (d *database) loadAccountsFromKeys(allKeys []string) ([]Account, error) {
	db, err := d.state.getOpenedDB()
	if err != nil {
		return nil, err
	}
	accs := make([]Account, 0)
	for _, k := range allKeys {
		err = db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(k))
			if err != nil {
				return err
			}
			acc := account{}
			err = item.Value(func(val []byte) (err error) {
				err = DecodeToStruct(&acc, val)
				return err
			})
			if err != nil {
				return err
			}
			accs = append(accs, Account{account: acc})
			return nil
		})
		if err != nil {
			continue
		}
	}
	return accs, err
}

func (d *database) deleteActiveAccount() error {
	db, err := d.state.getOpenedDB()
	if err != nil {
		return err
	}
	err = db.Update(func(txn *badger.Txn) error {
		err = txn.Delete([]byte(PathActiveAccount))
		return err
	})
	return err
}

func (d *database) activeAccount() (*Account, error) {
	db, err := d.state.getOpenedDB()
	if err != nil {
		return nil, err
	}
	var acc Account
	err = db.View(func(txn *badger.Txn) error {
		var item *badger.Item
		item, err = txn.Get([]byte(PathActiveAccount))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) (err error) {
			err = DecodeToStruct(&acc.account, val)
			if err != nil {
				return err
			}
			return nil
		})
		return err
	})
	return &acc, err
}

func (d *database) setActiveAccount(acc Account) error {
	db, err := d.state.getOpenedDB()
	if err != nil {
		return err
	}
	if acc.PathID() == "" {
		return ErrInvalidPathID
	}
	if acc.Type() == AccountTypeInvalid {
		return ErrInvalidAccount
	}
	val := EncodeToBytes(&acc.account)
	txn := db.NewTransaction(true)
	defer txn.Discard()
	err = txn.Set([]byte(PathActiveAccount), val)
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (d *database) saveAccount(pvtAcc *Account) error {
	db, err := d.state.getOpenedDB()
	if err != nil {
		return err
	}
	if pvtAcc.Type() == AccountTypeInvalid {
		return ErrInvalidAccount
	}
	key := fmt.Sprintf("%s/%s", PathAccounts, pvtAcc.PathID())
	val := EncodeToBytes(&pvtAcc.account)
	txn := db.NewTransaction(true)
	defer txn.Discard()
	err = txn.Set([]byte(key), val)
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}

// saveCollection
// if the PathType is PathTypeDocument then the /Collection.ID will be appended
// to the path
// if the PathType is PathTypeCollection then the last segment will be compared
// with the Collection.ID, on mismatch an error of ErrInvalidPath is returned
//func (d *database) saveCollection(collPath string, collection Collection) error {
//	dBase, err := d.state.getDB()
//	if err != nil {
//		return err
//	}
//	err = collPath.Validate()
//	if err != nil {
//		return err
//	}
//	if collPath.Type() != router.PathTypeCollection && collPath.Type() != router.PathTypeDocument {
//		return ErrInvalidPath
//	}
//	pathID, err := collPath.LastSegmentPathID()
//	if collPath.Type() == router.PathTypeCollection && collection.ID != pathID {
//		return ErrInvalidCollectionPath
//	}
//	if collPath.Type() == router.PathTypeDocument {
//		collPath = router.Path(fmt.Sprintf("%s/%s", collPath, collection.ID))
//	}
//	if err = collPath.Validate(); err != nil {
//		return err
//	}
//	txn := dBase.NewTransaction(true)
//	defer txn.Discard()
//	for _, doc := range collection.Values {
//		k := fmt.Sprintf("%s/%s", collPath, doc.ID)
//		err = txn.Set([]byte(k), doc.Value)
//		if err != nil {
//			return err
//		}
//	}
//	err = txn.Commit()
//	return err
//}

//// saveDocument
//// if the PathType is PathTypeCollection then the /Document.ID will be appended
//// to the path
//// if the PathType is PathTypeDocument then the last segment will be compared
//// with the Document.ID, on mismatch an error of ErrInvalidPath is returned
//func (d *DB) saveDocument(docPath router.Path, document router.Document) error {
//	dBase, err := d.state.getDB()
//	if err != nil {
//		return err
//	}
//	err = docPath.Validate()
//	if err != nil {
//		return err
//	}
//	if docPath.Type() != router.PathTypeCollection && docPath.Type() != router.PathTypeDocument {
//		return ErrInvalidPath
//	}
//	pathID, err := docPath.LastSegmentPathID()
//	if docPath.Type() == router.PathTypeDocument && document.ID != pathID {
//		return ErrInvalidDocumentPath
//	}
//	if docPath.Type() == router.PathTypeCollection {
//		docPath = router.Path(fmt.Sprintf("%s/%s", docPath, document.ID))
//	}
//	if err = docPath.Validate(); err != nil {
//		return err
//	}
//	txn := dBase.NewTransaction(true)
//	defer txn.Discard()
//	err = txn.Set([]byte(docPath), document.Value)
//	if err != nil {
//		return err
//	}
//	err = txn.Commit()
//	return err
//}
