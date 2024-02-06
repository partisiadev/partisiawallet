package evm

import "math/big"

type Chain struct {
	Name           string         `json:"name"`
	Chain          string         `json:"chain"`
	Icon           string         `json:"icon,omitempty"`
	IconsData      []IconData     `json:",omitempty"`
	RPC            []RPC          `json:"rpc"`
	Features       []Features     `json:"features,omitempty"`
	Faucets        []string       `json:"faucets"`
	NativeCurrency NativeCurrency `json:"nativeCurrency"`
	InfoURL        string         `json:"infoURL"`
	ShortName      string         `json:"shortName"`
	ChainID        big.Int        `json:"chainId"`
	NetworkID      big.Int        `json:"networkId"`
	Slip44         big.Int        `json:"slip44,omitempty"`
	Ens            Ens            `json:"ens,omitempty"`
	Explorers      []Explorer     `json:"explorers,omitempty"`
	Title          string         `json:"title,omitempty"`
	RedFlags       []string       `json:"redFlags,omitempty"`
	Parent         Parent         `json:"parent,omitempty"`
	Status         string         `json:"status,omitempty"`
}
