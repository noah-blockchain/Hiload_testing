package app

import (
	"github.com/noah-blockchain/Hiload_testing/internal/dao"
	"github.com/noah-blockchain/go-sdk/api"
	"os"
)

const (
	walletCount        = 1000
	countRandomLetters = 10
	maximumSendNoah    = 26000
)

var (
	seedFrom = os.Getenv("SEED_PHRASE")
)

type App interface {
	GetCountWallets() (uint64, error)
	CreateWallet(wallet Wallet) error
	CreateWallets() error
	Start() error
	UpdateWallets() error
}

type Repo interface {
	GetCountWallets() (uint64, error)
	CreateWallet(address, seedPhrase, mnemonic, privateKey, amount string, status bool) error
	SelectWallets() ([]dao.Wallet, error)
	SelectWalletsInterval(start, end uint64) ([]dao.Wallet, error)
	SelectWalletsAmount(amount uint64) ([]dao.Wallet, error)
}

type app struct {
	repo    Repo
	rl      RateLimiter
	nodeAPI *api.Api
}

func New(repo Repo, rl RateLimiter) App {
	a := &app{
		repo:    repo,
		rl:      rl,
		nodeAPI: api.NewApi(os.Getenv("NODE_API_URL")),
	}

	return a
}
