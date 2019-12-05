package app

import (
	"fmt"
	"github.com/noah-blockchain/go-sdk/transaction"
	"github.com/noah-blockchain/go-sdk/wallet"
	"math/big"
)

func (a app) createTokenTrx(name string) error {
	wallets, err := a.repo.SelectWalletsInterval(1127, 1131)
	if err != nil || len(wallets) == 0 {
		return err
	}

	trxs := make([]transaction.SignedTransaction, len(wallets))
	for i, w := range wallets {
		seed, _ := wallet.Seed(w.Mnemonic)
		walletFrom, err := wallet.NewWallet([]byte(seed))
		if err != nil {
			fmt.Println(err)
			continue
		}

		nonce, err := a.nodeAPI.Nonce(walletFrom.Address())
		if err != nil {
			fmt.Println(err)
			continue
		}

		tx := transaction.NewCreateCoinData()
		tx.SetConstantReserveRatio(100)
		amount, _ := big.NewInt(0).SetString("1000000000000000000", 10)
		tx.SetInitialAmount(amount)
		reserve, _ := big.NewInt(0).SetString("1000000000000000000000", 10)
		tx.SetInitialReserve(reserve)
		tx.SetName(name)
		tx.SetSymbol(name)
		signedTx, err := transaction.NewBuilder(transaction.TestNetChainID).NewTransaction(tx)
		if err != nil {
			fmt.Println(err)
			continue
		}

		finishedTx, err := signedTx.
			SetNonce(nonce).
			SetGasPrice(255).
			SetGasCoin("NOAH").
			SetPayload([]byte("TOKEN CREATED")).
			Sign(walletFrom.PrivateKey())
		if err != nil {
			fmt.Println(err)
			continue
		}

		trxs[i] = finishedTx
	}

	for _, trx := range trxs {
		//go func(signTrx transaction.SignedTransaction) {
		fmt.Println("START TRX")
		res, err := a.nodeAPI.SendTransaction(trx)
		if err != nil {
			fmt.Println(err)
			return err
		}
		a, _ := trx.SenderAddress()
		fmt.Println(a)
		fmt.Println(res.Hash)
		fmt.Println(res.Code)
		fmt.Println(res.Log)
		//}(trx)
	}

	return nil
}
