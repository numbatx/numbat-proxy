package address

import "github.com/numbatx/numbat-proxy/data"

// FacadeHandler interface defines methods that can be used from `numbatProxyFacade` context variable
type FacadeHandler interface {
	GetAccount(address string) (*data.Account, error)
}
