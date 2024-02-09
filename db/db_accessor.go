package db

import (
	"fmt"
	"gioui.org/app"
	"github.com/dgraph-io/badger/v4"
	"github.com/partisiadev/partisiawallet/config"
	"github.com/partisiadev/partisiawallet/log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const MaxNumOfPasswdChars = 32
const passwdPadCharacter = "0"

type Accessor struct {
	password      string
	passwordMutex sync.RWMutex
	dataBase      *badger.DB
	dataBaseMutex sync.RWMutex
}

var dbState = &Accessor{}

func accessorInstance() *Accessor {
	return dbState
}

func (s *Accessor) setPassword(passwd string) {
	s.passwordMutex.Lock()
	defer s.passwordMutex.Unlock()
	s.password = passwd
}

func (s *Accessor) getPassword() string {
	s.passwordMutex.RLock()
	defer s.passwordMutex.RUnlock()
	return s.password
}

func (s *Accessor) setDB(db *badger.DB) {
	s.dataBaseMutex.Lock()
	defer s.dataBaseMutex.Unlock()
	s.dataBase = db
}

func (s *Accessor) getDB() *badger.DB {
	s.dataBaseMutex.RLock()
	defer s.dataBaseMutex.RUnlock()
	return s.dataBase
}

func (s *Accessor) getOpenedDB() (*badger.DB, error) {
	db := s.getDB()
	if db == nil {
		return nil, ErrDBIsClosed
	}
	if db.IsClosed() {
		s.setDB(nil)
		return nil, ErrDBIsClosed
	}
	return db, nil
}

func (s *Accessor) OpenDB(passwd string) error {
	db := s.getDB()
	if db != nil || s.IsDBOpen() {
		return ErrDBIsAlreadyOpen
	}
	if len(passwd) == 0 {
		return ErrPasswdCannotBeEmpty
	}
	if len([]byte(passwd)) > MaxNumOfPasswdChars {
		return fmt.Errorf("password should be less than %d characters", MaxNumOfPasswdChars)
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
	origPasswd := passwd
	passwd = leftPad + passwd + rightPad
	dbPath := s.DatabasePath()
	options := badger.DefaultOptions(dbPath)
	options.EncryptionKey = []byte(passwd)
	options.IndexCacheSize = 100 << 20
	dataBase, err := badger.Open(options)
	if err != nil {
		return err
	}
	s.setDB(dataBase)
	s.setPassword(origPasswd)
	return nil
}

func (s *Accessor) DatabasePath() string {
	dirPath, err := app.DataDir()
	if err != nil {
		log.Logger().Fatal(err)
	}
	dbPath := filepath.Join(dirPath, config.AppConfigPathDirName, config.AppWalletDBDirName)
	return dbPath
}

func (s *Accessor) DatabaseExists() bool {
	dbPath := s.DatabasePath()
	_, err := os.Stat(dbPath)
	if err == nil {
		res, err := os.ReadDir(dbPath)
		// If database directory is empty, then
		if len(res) == 0 && err == nil {
			err = os.RemoveAll(dbPath)
			return false
		}
		return true
	}
	return err == nil
}

// VerifyPassword
// returns nil if password is correct else may
// return ErrPasswordMismatch or ErrPasswordInvalid or ErrPasswordNotSet
func (s *Accessor) VerifyPassword(passwd string) error {
	if strings.TrimSpace(passwd) == "" {
		return ErrPasswordInvalid
	}
	if s.getPassword() == "" {
		return ErrPasswordNotSet
	}
	if s.getPassword() != passwd {
		return ErrPasswordMismatch
	}
	return nil
}

func (s *Accessor) closeDB() error {
	db := s.getDB()
	if db == nil {
		_ = db.Close()
	}
	s.setDB(nil)
	return nil
}

func (s *Accessor) IsDBOpen() bool {
	db := s.getDB()
	if db == nil {
		return false
	}
	return !db.IsClosed()
}
