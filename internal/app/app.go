package app

import (
	"os"
)

const (
	walletCount        = 5
	countRandomLetters = 10
	maximumSendNoah    = 1000
)

var (
	seedFrom = os.Getenv("SEED_PHRASE")
)

type App interface {
	GetCountWallets() (uint64, error)
	CreateWallet(wallet Wallet) error
	CheckAndCreateWallets() error
}

type Repo interface {
	GetCountWallets() (uint64, error)
	CreateWallet(address, seedPhrase, mnemonic, privateKey, amount string, status bool) error
}

type app struct {
	repo Repo
}

func New(repo Repo) App {
	a := &app{
		repo: repo,
	}

	return a
}
