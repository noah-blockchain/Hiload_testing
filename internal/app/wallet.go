package app

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/noah-blockchain/Hiload_testing/internal/utils"
	"github.com/noah-blockchain/go-sdk/transaction"
	"github.com/noah-blockchain/go-sdk/wallet"
)

var walletNonce uint64

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

func (a *app) CreateWallets() error {
	var mustBeCreated = walletCount
	i := 0
	for {
		count := 0
		if mustBeCreated <= 100 {
			count = mustBeCreated
			mustBeCreated = 0
		} else {
			count = 100
			mustBeCreated -= count
		}
		wg := sync.WaitGroup{}
		wg.Add(count)
		wallets := make([]Wallet, count)
		for i := 0; i < count; i++ {
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

		multiSendList := make([]MultiSendItem, count)
		s := i * count
		e := s + count
		for i, w := range wallets[s:e] {
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
			continue
		}

		i++
		if mustBeCreated == 0 {
			break
		}

		time.Sleep(time.Second)
	}

	return nil
}

func (a *app) UpdateWallets() error {
	wallets, err := a.repo.SelectWallets()
	if err != nil {
		return err
	}
	wallets = wallets[0:384]
	var mustBeUpdated = len(wallets)
	i := 0
	for {
		count := 0
		if mustBeUpdated <= 100 {
			count = mustBeUpdated
			mustBeUpdated = 0
		} else {
			count = 100
			mustBeUpdated -= count
		}

		multiSendList := make([]MultiSendItem, count)
		s := i * count
		e := s + count
		for i, w := range wallets[s:e] {
			fmt.Println(fmt.Sprintf("Send %d NOAH to %s", maximumSendNoah, w.Address))
			valueQnoah := utils.NoahToQNoah(big.NewInt(maximumSendNoah))
			multiSendList[i] = MultiSendItem{
				Coin:  "NOAH",
				To:    w.Address,
				Value: valueQnoah,
				wallet: Wallet{
					Address:    w.Address,
					PrivateKey: w.PrivateKey,
					SeedPhrase: w.SeedPhrase,
					Mnemonic:   w.Mnemonic,
				},
			}
		}

		if err := a.sendMultiListTrx(multiSendList, utils.String(countRandomLetters)); err != nil {
			fmt.Println(err)
			continue
		}

		i++
		if mustBeUpdated == 0 {
			break
		}

		time.Sleep(time.Second)
	}

	return nil
}

func (a *app) sendMultiListTrx(list []MultiSendItem, payload string) error {
	seed, _ := wallet.Seed(seedFrom)
	walletFrom, err := wallet.NewWallet([]byte(seed))
	if err != nil {
		return err
	}

	m := sync.RWMutex{}
	for {
		nonce, err := a.GetNodeURL().Nonce(walletFrom.Address())
		if err != nil {
			return err
		}

		if nonce > walletNonce {
			m.Lock()
			walletNonce = nonce
			m.Unlock()
			break
		}

		time.Sleep(time.Second)
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
		SetNonce(walletNonce).
		SetGasPrice(255).
		SetGasCoin("NOAH").
		SetPayload([]byte(payload)).
		Sign(walletFrom.PrivateKey())
	if err != nil {
		return err
	}

	_, err = a.GetNodeURL().SendTransaction(finishedTx)
	if err != nil {
		return err
	}

	for _, d := range list {
		err = a.repo.CreateWallet(
			d.wallet.Address,
			d.wallet.SeedPhrase,
			d.wallet.Mnemonic,
			d.wallet.PrivateKey,
			utils.QNoahStr2Noah(d.Value.String()),
			true,
		)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
