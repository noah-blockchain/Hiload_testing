package dal

import (
	"github.com/noah-blockchain/Hiload_testing/internal/dao"
)

func (r repo) GetCountWallets() (uint64, error) {
	var count uint64
	rows, err := r.db.Query("SELECT COUNT(*) as count FROM wallets")
	if err != nil {
		return 0, nil
	}

	for rows.Next() {
		err = rows.Scan(&count)
	}
	if err != nil {
		return 0, nil
	}

	return count, nil
}

func (r repo) CreateWallet(address, seedPhrase, mnemonic, privateKey, amount string, status bool) error {
	walletDao := dao.Wallet{
		Address:    address,
		SeedPhrase: seedPhrase,
		Mnemonic:   mnemonic,
		PrivateKey: privateKey,
		Amount:     amount,
		Status:     status,
	}

	createWalletSQL := `INSERT INTO wallets (address, seed_phrase, mnemonic, private_key, amount, status) VALUES (:address, :seed_phrase, :mnemonic, :private_key, :amount, :status)`
	tx := r.db.MustBegin()
	_, err := tx.NamedExec(createWalletSQL, &walletDao)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}