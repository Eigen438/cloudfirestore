package mock

import (
	"context"

	"github.com/Eigen438/cloudfirestore"
	"github.com/stretchr/testify/mock"
)

type innerTran struct {
	mock *mock.Mock
	tran cloudfirestore.Transaction
}

func (i *innerTran) Create(ctx context.Context, data cloudfirestore.KeyGenerator) error {
	args := i.mock.Called(ctx, data)
	err := args.Error(0)
	if err != nil {
		return err
	}
	return i.tran.Create(ctx, data)
}

func (i *innerTran) Set(ctx context.Context, data cloudfirestore.KeyGenerator) error {
	args := i.mock.Called(ctx, data)
	err := args.Error(0)
	if err != nil {
		return err
	}
	return i.tran.Set(ctx, data)
}

func (i *innerTran) Get(ctx context.Context, data cloudfirestore.KeyGenerator) error {
	args := i.mock.Called(ctx, data)
	err := args.Error(0)
	if err != nil {
		return err
	}
	return i.tran.Get(ctx, data)
}

func (i *innerTran) Delete(ctx context.Context, data cloudfirestore.KeyGenerator) error {
	args := i.mock.Called(ctx, data)
	err := args.Error(0)
	if err != nil {
		return err
	}
	return i.tran.Delete(ctx, data)
}
