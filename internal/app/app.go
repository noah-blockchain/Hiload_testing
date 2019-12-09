package app

import (
	"fmt"
	"github.com/noah-blockchain/Hiload_testing/internal/env"
	"math/rand"
	"os"
	"strings"

	"github.com/noah-blockchain/Hiload_testing/internal/dao"
	"github.com/noah-blockchain/go-sdk/api"
)

const (
	countRandomLetters = 10
	maximumSendNoah    = 10
)

var (
	seedFrom    = os.Getenv("SEED_PHRASE")
	walletCount = env.GetEnvAsInt("WALLET_COUNT", 500000)
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
	DisableWallet(address string) error
	GetOneWallet() (*dao.Wallet, error)
}

type app struct {
	repo     Repo
	rl       RateLimiter
	nodeAPIs []*api.Api
}

func New(repo Repo, rl RateLimiter) App {
	apis := strings.Split(os.Getenv("NODE_API_URLS"), ",")

	nodeAPIs := make([]*api.Api, len(apis))
	for i, a := range apis {
		nodeAPIs[i] = api.NewApi(a)
		fmt.Println("Node", i, "URL", a)
	}

	a := &app{
		repo:     repo,
		rl:       rl,
		nodeAPIs: nodeAPIs,
	}

	return a
}

func (a app) GetNodeURL() *api.Api {
	return a.nodeAPIs[rand.Int31n(int32(len(a.nodeAPIs)))]
}