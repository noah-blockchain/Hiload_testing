package app

import (
	"fmt"
	"github.com/noah-blockchain/Hiload_testing/internal/utils"
	"github.com/noah-blockchain/go-sdk/api"
	"github.com/noah-blockchain/go-sdk/transaction"
	"github.com/noah-blockchain/go-sdk/wallet"
	"math/big"
	"os"
	"sync"
)

type Wallet struct {
	Address    string
	PrivateKey string
	SeedPhrase string
	Mnemonic   string
}

func (a app) GetCountWallets() (uint64, error) {
	panic("implement me")
}

func (a *app) CreateWallet(wallet Wallet) error {
	panic("implement me")
}

func (a *app) CheckAndCreateWallets() error {
	count, err := a.repo.GetCountWallets()
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(walletCount)
	wallets := make([]Wallet, walletCount)
	for i := 0; i < walletCount; i++ {
		go func(i int) {
			walletData, _ := wallet.Create()
			wallets[i] = Wallet{
				Address:    walletData.Address,
				PrivateKey: walletData.PrivateKey,
				Mnemonic:   walletData.Mnemonic,
				SeedPhrase: walletData.Seed,
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	multiSendList := make([]MultiSendItem, len(wallets))
	for i, w := range wallets {
		fmt.Println(fmt.Sprintf("Send %d NOAH to %s", maximumSendNoah, w.Address))
		valueQnoah := utils.NoahToQNoah(big.NewInt(maximumSendNoah))
		multiSendList[i] = MultiSendItem{
			Coin:   "NOAH",
			To:     w.Address,
			Value:  valueQnoah,
			wallet: w,
		}
	}

	if err := a.sendMultiListTrx(multiSendList, utils.String(countRandomLetters)); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (a *app) sendMultiListTrx(list []MultiSendItem, payload string) error {
	nodeAPI := api.NewApi(os.Getenv("NODE_API_URL"))

	seed, _ := wallet.Seed(seedFrom)
	walletFrom, err := wallet.NewWallet([]byte(seed))
	if err != nil {
		return err
	}

	nonce, err := nodeAPI.Nonce(walletFrom.Address())
	if err != nil {
		return err
	}

	tx := transaction.NewMultisendData()
	for _, d := range list {
		tx.AddItem(
			*transaction.NewMultisendDataItem().
				SetCoin(d.Coin).
				SetValue(d.Value).
				MustSetTo(d.To),
		)
	}

	signedTx, err := transaction.NewBuilder(transaction.TestNetChainID).NewTransaction(tx)
	if err != nil {
		return err
	}

	//comission := maximumSendNoah * 0.1 - 0.02
	//payment := maximumSendNoah - comission - 0.01
	finishedTx, err := signedTx.
		SetNonce(nonce).
		SetGasPrice(255).
		SetGasCoin("NOAH").
		SetPayload([]byte(payload)).
		Sign(walletFrom.PrivateKey())
	if err != nil {
		return err
	}

	_, err = nodeAPI.SendTransaction(finishedTx)
	if err != nil {
		return err
	}

	for _, d := range list {
		err = a.repo.CreateWallet(
			d.wallet.Address,
			d.wallet.SeedPhrase,
			d.wallet.Mnemonic,
			d.wallet.PrivateKey,
			d.Value.String(),
			true,
		)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
