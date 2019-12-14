package app

import (
	"fmt"
	"github.com/noah-blockchain/Hiload_testing/internal/utils"
	"github.com/noah-blockchain/go-sdk/api"
	"github.com/noah-blockchain/go-sdk/transaction"
	"github.com/noah-blockchain/go-sdk/wallet"
	"math/big"
	"sync"
	"time"
)

var nonceWallets sync.Map

func (a app) Start() error {
	//if err := a.createTokenTrx(strings.ToUpper(utils.String(8))); err != nil {
	//	fmt.Println(err)
	//}

	if err := a.createSendTrx(); err != nil {
		fmt.Println(err)
	}

	return nil
}

func (a app) createSendTrx() error {
	count := 0
	for {
		for i := 0; i < a.rl.Freq; i++ {
			walletFrom, _ := a.repo.GetOneWallet()
			walletTo, _ := a.repo.GetOneWallet()

			if walletFrom == nil || walletTo == nil {
				continue
			}

			go func(from, to string) {
				a.sendNoahFromTo(from, to)
			}(walletFrom.Mnemonic, walletTo.Address)
		}
		count++
		fmt.Println(fmt.Sprintf("Maked %d trx", a.rl.Freq*count))
		fmt.Println("Go to sleep", a.rl.Per.String())
		time.Sleep(a.rl.Per)
	}
}

func (a app) sendNoahFromTo(fromMnemonic string, toAddress string) {
	seed, _ := wallet.Seed(fromMnemonic)
	walletFrom, err := wallet.NewWallet([]byte(seed))
	if err != nil {
		fmt.Println(err)
		return
	}

	nonce, err := a.GetNodeURL().Nonce(walletFrom.Address())
	v, ok := nonceWallets.Load(walletFrom.Address())
	if ok {
		value := v.(uint64)
		for {
			if nonce > value {
				nonceWallets.Store(walletFrom.Address(), nonce)
				break
			}

			nonce, err = a.GetNodeURL().Nonce(walletFrom.Address())
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	tx := transaction.NewSendData()
	tx.SetValue(utils.NoahToQNoah(big.NewInt(1)))
	tx.SetCoin("NOAH")
	_, _ = tx.SetTo(toAddress)

	signedTx, err := transaction.NewBuilder(transaction.TestNetChainID).NewTransaction(tx)
	if err != nil {
		fmt.Println(err)
		return
	}

	finishedTx, err := signedTx.
		SetNonce(nonce).
		SetGasPrice(255).
		SetGasCoin("NOAH").
		SetPayload([]byte("SEND TRX")).
		Sign(walletFrom.PrivateKey())
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := a.GetNodeURL().SendTransaction(finishedTx)
	if err != nil {
		v, ok := err.(*api.TxError)
		if ok && v.Code == 412 {
			_ = a.repo.DisableWallet(walletFrom.Address())
		}
		fmt.Println(err)
		return
	}
	fmt.Println("Success", res.Hash)
}
