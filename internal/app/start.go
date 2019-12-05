package app

import (
	"fmt"
	"github.com/noah-blockchain/Hiload_testing/internal/utils"
	"github.com/noah-blockchain/go-sdk/transaction"
	"github.com/noah-blockchain/go-sdk/wallet"
	"math/big"
	"math/rand"
	"time"
)

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
	wallets, err := a.repo.SelectWallets()
	if err != nil || len(wallets) == 0 {
		return err
	}
	//began, count := time.Now(), uint64(0)
	count := 0

	for {
		fmt.Println("Lets go")
		//elapsed := time.Since(began)
		//_, stop := a.rl.Pace(elapsed, count)
		//if stop {
		//	break
		//}
		//fmt.Println("Sleep", wait)
		//time.Sleep(wait)
		//for i := 0; i < a.rl.Freq; i++ {
		for i := 0; i < 1000; i++ {
			walletFrom := wallets[rand.Intn(len(wallets))]
			walletTo := wallets[rand.Intn(len(wallets))]
			go func(from, to string) {
				a.sendNoahFromTo(from, to)
			}(walletFrom.Mnemonic, walletTo.Address)
			//a.sendNoahFromTo(walletFrom.Mnemonic, walletTo.Address)
		}
		count++
		//fmt.Println(fmt.Sprintf("Maked %d trx", 150*count))
		fmt.Println(fmt.Sprintf("Maked %d trx", 250*count))
		//fmt.Println("Go to sleep 1 minute")
		fmt.Println("Go to sleep 1 second")
		//time.Sleep(time.Minute)
		time.Sleep(5 * time.Second)
	}

}

func (a app) sendNoahFromTo(fromMnemonic string, toAddress string) {
	seed, _ := wallet.Seed(fromMnemonic)
	walletFrom, err := wallet.NewWallet([]byte(seed))
	if err != nil {
		fmt.Println(err)
		return
	}

	nonce, err := a.nodeAPI.Nonce(walletFrom.Address())
	if err != nil {
		fmt.Println(err)
		return
	}

	tx := transaction.NewSendData()
	tx.SetValue(utils.NoahToQNoah(big.NewInt(10)))
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

	fmt.Println("START TRX")
	res, err := a.nodeAPI.SendTransaction(finishedTx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Hash)
}
