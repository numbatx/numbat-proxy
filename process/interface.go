package process

import (
	"github.com/numbatx/numbat-proxy/config"
	"github.com/numbatx/numbat-proxy/data"
)

// Processor defines what a processor should be able to do
type Processor interface {
	ApplyConfig(cfg *config.Config) error
	GetObservers(shardId uint32) ([]*data.Observer, error)
	ComputeShardId(addressBuff []byte) (uint32, error)
	CallGetRestEndPoint(address string, path string, value interface{}) error
	CallPostRestEndPoint(address string, path string, data interface{}, response interface{}) error
}
