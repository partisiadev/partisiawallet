package db

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"strings"
)

func autoCreateECDSAAccount() (*Account, error) {
	var acc Account
	ecdsaPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	publicKey := ecdsaPrivateKey.PublicKey
	pvtKeyStr := hex.EncodeToString(ecdsaPrivateKey.D.Bytes())
	publicKeyStr := hex.EncodeToString(crypto.FromECDSAPub(&publicKey))
	ethAddress := hex.EncodeToString(crypto.PubkeyToAddress(publicKey).Bytes())
	acc.account = account{
		PrivateKey: pvtKeyStr,
		PublicKey:  publicKeyStr,
		EthAddress: ethAddress,
	}
	return &acc, err
}

func importECDSAAccount(pvtKeyOrMnemonic string) (*Account, error) {
	var acc Account
	var err error
	if strings.TrimSpace(pvtKeyOrMnemonic) == "" {
		err = errors.New("private key is empty")
		return nil, err
	}
	var pvtKey *ecdsa.PrivateKey
	if len(strings.Split(pvtKeyOrMnemonic, " ")) > 11 {
		acc.account.Mnemonics = pvtKeyOrMnemonic
		// Generate a Bip32 HD wallet for the mnemonic and a user supplied passphrase
		seed := bip39.NewSeed(pvtKeyOrMnemonic, "")
		masterPrivateKey, _ := bip32.NewMasterKey(seed)
		// masterPublicKey := masterPrivateKey.PublicKey()
		// BIP44 derivation path format: m / purpose' / coin_type' / account' / change / address_index
		// Example: m/44'/60'/0'/0/0
		purposeKey, _ := masterPrivateKey.NewChildKey(bip32.FirstHardenedChild + 44)
		coinTypeKey, _ := purposeKey.NewChildKey(bip32.FirstHardenedChild + 60)
		accountKey, _ := coinTypeKey.NewChildKey(bip32.FirstHardenedChild)
		changeKey, _ := accountKey.NewChildKey(0)
		addressKey, _ := changeKey.NewChildKey(0)
		pvtKey, err = crypto.ToECDSA(addressKey.Key)
		if err != nil {
			return nil, err
		}
		pvtKeyOrMnemonic = fmt.Sprintf("%x", pvtKey.D)
	} else {
		pvtKey, err = crypto.HexToECDSA(pvtKeyOrMnemonic)
		if err != nil {
			return nil, err
		}
	}
	publicKeyStr := fmt.Sprintf("%x", crypto.CompressPubkey(pvtKey.Public().(*ecdsa.PublicKey)))
	ethAddress := hex.EncodeToString(crypto.PubkeyToAddress(pvtKey.PublicKey).Bytes())
	acc.account.PrivateKey = pvtKeyOrMnemonic
	acc.account.PublicKey = publicKeyStr
	acc.account.EthAddress = ethAddress
	return &acc, nil
}
