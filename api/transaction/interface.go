package transaction

import (
	"math/big"
)

// FacadeHandler interface defines methods that can be used from `numbatProxyFacade` context variable
type FacadeHandler interface {
	SendTransaction(nonce uint64, sender string, receiver string, value *big.Int, code string, signature []byte) (string, error)
}
