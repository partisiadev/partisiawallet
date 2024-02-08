package db

import (
	"errors"
	"github.com/partisiadev/partisiawallet/db/internal"
)

var (
	ErrInvalidCollectionPath = errors.New("collection path is not valid")
	ErrInvalidDocumentPath   = errors.New("document path is not valid")
	ErrInvalidPath           = errors.New("invalid path")
)

type State interface {
	OpenDB(passwd string) error
	DatabaseExists() bool
	VerifyPassword(passwd string) error
	IsDBOpen() bool
}

type DB struct {
	state *internal.State
}

var dbInstance = &DB{}
var _ State = dbInstance.state

func Instance() *DB {
	if dbInstance.state == nil {
		dbInstance.state = internal.DBState()
	}
	return dbInstance
}

func (d *DB) State() State {
	return d.state
}

//
//// saveCollection
//// if the PathType is PathTypeDocument then the /Collection.ID will be appended
//// to the path
//// if the PathType is PathTypeCollection then the last segment will be compared
//// with the Collection.ID, on mismatch an error of ErrInvalidPath is returned
//func (d *DB) saveCollection(collPath router.Path, collection router.Collection) error {
//	dBase, err := d.state.GetDB()
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
//
//// saveDocument
//// if the PathType is PathTypeCollection then the /Document.ID will be appended
//// to the path
//// if the PathType is PathTypeDocument then the last segment will be compared
//// with the Document.ID, on mismatch an error of ErrInvalidPath is returned
//func (d *DB) saveDocument(docPath router.Path, document router.Document) error {
//	dBase, err := d.state.GetDB()
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
