package facade

import (
	"math/big"

	"github.com/numbatx/numbat-proxy/data"
)

// NumbatProxyFacade implements the facade used in api calls
type NumbatProxyFacade struct {
	accountProc AccountProcessor
	txProc      TransactionProcessor
}

// NewNumbatProxyFacade creates a new NumbatProxyFacade instance
func NewNumbatProxyFacade(
	accountProc AccountProcessor,
	txProc TransactionProcessor,
) (*NumbatProxyFacade, error) {

	if accountProc == nil {
		return nil, ErrNilAccountProcessor
	}
	if txProc == nil {
		return nil, ErrNilTransactionProcessor
	}

	return &NumbatProxyFacade{
		accountProc: accountProc,
		txProc:      txProc,
	}, nil
}

// GetAccount returns an account based on the input address
func (epf *NumbatProxyFacade) GetAccount(address string) (*data.Account, error) {
	return epf.accountProc.GetAccount(address)
}

// SendTransaction should sends the transaction to the correct observer
func (epf *NumbatProxyFacade) SendTransaction(
	nonce uint64,
	sender string,
	receiver string,
	value *big.Int,
	code string,
	signature []byte,
) (string, error) {

	return epf.txProc.SendTransaction(nonce, sender, receiver, value, code, signature)
}
