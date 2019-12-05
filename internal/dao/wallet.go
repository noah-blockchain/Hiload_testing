package dao

type Wallet struct {
	Id         int    `db:"id"`
	Address    string `db:"address"`
	PrivateKey string `db:"private_key"`
	SeedPhrase string `db:"seed_phrase"`
	Mnemonic   string `db:"mnemonic"`
	Amount     string `db:"amount" sql:"type:numeric(70)"`
	Status     bool   `db:"status"`
}
