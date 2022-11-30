package assets

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LogTransfer struct {
	TokenId *big.Int
}

type Token struct {
	Token string  `json:"token"`
	Ids   []int64 `json:"ids"`
}

type Assets struct {
	Tokens []*Token `json:"assets"`
}

var contracts = []common.Address{
	common.HexToAddress("0x61B91a780945971b07ba3898A8E0Dc8201dB46b3"),
	common.HexToAddress("0xd89148d8dEFc7E8942cA4b16DdBE0E2f6485a4c8"),
	common.HexToAddress("0x9D7f9672060EED641ebc1b22443132eDf4967D91"),
}
