package db

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
)

func autoCreateEcdsaAccount() (*Account, error) {
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
