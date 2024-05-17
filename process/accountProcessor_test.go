package process_test

import (
	"errors"
	"testing"

	"github.com/numbatx/numbat-proxy/data"
	"github.com/numbatx/numbat-proxy/process"
	"github.com/numbatx/numbat-proxy/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewAccountProcessor_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	ap, err := process.NewAccountProcessor(nil)

	assert.Nil(t, ap)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewAccountProcessor_WithCoreProcessorShouldWork(t *testing.T) {
	t.Parallel()

	ap, err := process.NewAccountProcessor(&mock.ProcessorStub{})

	assert.NotNil(t, ap)
	assert.Nil(t, err)
}

//------- GetAccount

func TestAccountProcessor_GetAccountInvalidHexAdressShouldErr(t *testing.T) {
	t.Parallel()

	ap, _ := process.NewAccountProcessor(&mock.ProcessorStub{})
	accnt, err := ap.GetAccount("invalid hex number")

	assert.Nil(t, accnt)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
}

func TestAccountProcessor_GetAccountComputeShardIdFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	ap, _ := process.NewAccountProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, errExpected
		},
	})
	address := "DEADBEEF"
	accnt, err := ap.GetAccount(address)

	assert.Nil(t, accnt)
	assert.Equal(t, errExpected, err)
}

func TestAccountProcessor_GetAccountGetObserversFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	ap, _ := process.NewAccountProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return nil, errExpected
		},
	})
	address := "DEADBEEF"
	accnt, err := ap.GetAccount(address)

	assert.Nil(t, accnt)
	assert.Equal(t, errExpected, err)
}

func TestAccountProcessor_GetAccountSendingFailsOnAllObserversShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	ap, _ := process.NewAccountProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return []*data.Observer{
				{Address: "adress1", ShardId: 0},
				{Address: "adress2", ShardId: 0},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			return errExpected
		},
	})
	address := "DEADBEEF"
	accnt, err := ap.GetAccount(address)

	assert.Nil(t, accnt)
	assert.Equal(t, process.ErrSendingRequest, err)
}

func TestAccountProcessor_GetAccountSendingFailsOnFirstObserverShouldStillSend(t *testing.T) {
	t.Parallel()

	addressFail := "address1"
	errExpected := errors.New("expected error")
	respondedAccount := &data.ResponseAccount{
		AccountData: data.Account{
			Address: "an address",
		},
	}
	ap, _ := process.NewAccountProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return []*data.Observer{
				{Address: addressFail, ShardId: 0},
				{Address: "adress2", ShardId: 0},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			if address == addressFail {
				return errExpected
			}

			valRespond := value.(*data.ResponseAccount)
			valRespond.AccountData = respondedAccount.AccountData
			return nil
		},
	})
	address := "DEADBEEF"
	accnt, err := ap.GetAccount(address)

	assert.Equal(t, &respondedAccount.AccountData, accnt)
	assert.Nil(t, err)
}
