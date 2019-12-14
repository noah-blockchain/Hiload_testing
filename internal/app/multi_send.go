package app

import "math/big"

type MultiSendItem struct {
	Coin   string
	To     string
	Value  *big.Int
	wallet Wallet
}
