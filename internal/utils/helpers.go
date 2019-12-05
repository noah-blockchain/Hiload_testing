package utils

import "math/big"

func NoahToQNoah(noah *big.Int) *big.Int {
	p := big.NewInt(10)
	p.Exp(p, big.NewInt(18), nil)
	p.Mul(p, noah)

	return p
}

var qNoahInNoah = big.NewFloat(1000000000000000000)

func QNoahStr2Noah(value string) string {
	if value == "" {
		return "0"
	}

	floatValue, _ := new(big.Float).SetPrec(500).SetString(value)
	return new(big.Float).SetPrec(500).Quo(floatValue, qNoahInNoah).Text('f', 18)
}
