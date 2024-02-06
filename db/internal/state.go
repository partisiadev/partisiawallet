package internal

import (
	"errors"
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

var (
	ErrDBIsClosed          = errors.New("database is closed")
	ErrPasswordMismatch    = errors.New("password mismatch")
	ErrPasswordInvalid     = errors.New("password invalid")
	ErrPasswordNotSet      = errors.New("password not set")
	ErrPasswdCannotBeEmpty = errors.New("password cannot be empty")
)

type State struct {
	password      string
	passwordMutex sync.RWMutex
	dataBase      *badger.DB
	dataBaseMutex sync.RWMutex
}

var dbState = &State{}

func DBState() *State {
	return dbState
}

func (s *State) setPassword(passwd string) {
	s.passwordMutex.Lock()
	defer s.passwordMutex.Unlock()
	s.password = passwd
}

func (s *State) getPassword() string {
	s.passwordMutex.RLock()
	defer s.passwordMutex.RUnlock()
	return s.password
}

func (s *State) setDB(db *badger.DB) {
	s.dataBaseMutex.Lock()
	defer s.dataBaseMutex.Unlock()
	s.dataBase = db
}

func (s *State) GetDB() (*badger.DB, error) {
	s.dataBaseMutex.RLock()
	defer s.dataBaseMutex.RUnlock()
	if s.dataBase == nil || s.dataBase.IsClosed() {
		return nil, ErrDBIsClosed
	}
	return s.dataBase, nil
}

func (s *State) OpenDB(passwd string) error {
	origPasswd := passwd
	_, err := s.GetDB()
	if err != nil {
		return err
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

func (s *State) DatabasePath() string {
	dirPath, err := app.DataDir()
	if err != nil {
		log.Logger().Fatal(err)
	}
	dbPath := filepath.Join(dirPath, config.AppConfigPathDirName, config.AppWalletDBDirName)
	return dbPath
}

func (s *State) DatabaseExists() bool {
	dbPath := s.DatabasePath()
	_, err := os.Stat(dbPath)
	return err == nil
}

// VerifyPassword
// returns nil if password is correct else may
// return ErrPasswordMismatch or ErrPasswordInvalid or ErrPasswordNotSet
func (s *State) VerifyPassword(passwd string) error {
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

func (s *State) closeDB() error {
	db, err := s.GetDB()
	if err == nil {
		_ = db.Close()
	}
	s.setDB(nil)
	return nil
}

func (s *State) IsDBOpen() bool {
	_, err := s.GetDB()
	return err == nil
}
