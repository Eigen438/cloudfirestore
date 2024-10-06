// MIT License
//
// Copyright (c) 2024 Eigen
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package mock

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/Eigen438/cloudfirestore"
	"github.com/stretchr/testify/mock"
)

type inner struct {
	mock   *mock.Mock
	client cloudfirestore.CloudFirestore
}

func New(m *mock.Mock, firestoreService cloudfirestore.CloudFirestore) cloudfirestore.CloudFirestore {
	return &inner{
		mock:   m,
		client: firestoreService,
	}
}

func (i *inner) Create(ctx context.Context, data cloudfirestore.KeyGenerator) error {
	args := i.mock.Called(ctx, data)
	if err := args.Error(0); err != nil {
		return err
	}
	return i.client.Create(ctx, data)
}

func (i *inner) Delete(ctx context.Context, data cloudfirestore.KeyGenerator) error {
	args := i.mock.Called(ctx, data)
	if err := args.Error(0); err != nil {
		return err
	}
	return i.client.Delete(ctx, data)
}

func (i *inner) Get(ctx context.Context, data cloudfirestore.KeyGenerator) error {
	args := i.mock.Called(ctx, data)
	if err := args.Error(0); err != nil {
		return err
	}
	return i.client.Get(ctx, data)
}

func (i *inner) Set(ctx context.Context, data cloudfirestore.KeyGenerator) error {
	args := i.mock.Called(ctx, data)
	if err := args.Error(0); err != nil {
		return err
	}
	return i.client.Set(ctx, data)
}

func (i *inner) RunTransaction(ctx context.Context, f func(context.Context, cloudfirestore.Transaction) error) error {
	args := i.mock.Called(ctx, f)
	if err := args.Error(0); err != nil {
		return err
	}
	return i.client.RunTransaction(ctx, func(_ctx context.Context, tran cloudfirestore.Transaction) error {
		tx := &innerTran{
			mock: i.mock,
			tran: tran,
		}
		return f(_ctx, tx)
	})
}

func (i *inner) Collection(collectionName string) firestore.Query {
	i.mock.Called(collectionName)
	return i.client.Collection(collectionName)
}

func (i *inner) CollectionGroup(collectionName string) firestore.Query {
	i.mock.Called(collectionName)
	return i.client.CollectionGroup(collectionName)
}

func (i *inner) Sequence(ctx context.Context, q firestore.Query, f func(ctx context.Context, snapshot *firestore.DocumentSnapshot) error) (int, error) {
	args := i.mock.Called(ctx, q, f)
	err := args.Error(1)
	if err != nil {
		return args.Int(0), err
	}
	return i.client.Sequence(ctx, q, f)
}

func (i *inner) Run(ctx context.Context, q firestore.Query, concurrency int, f func(ctx context.Context, snapshot *firestore.DocumentSnapshot) error) (int, error) {
	args := i.mock.Called(ctx, q, concurrency, f)
	err := args.Error(1)
	if err != nil {
		return args.Int(0), err
	}
	return i.client.Run(ctx, q, concurrency, f)
}

func (i *inner) DeleteWithQuery(ctx context.Context, q firestore.Query, concurrency int) (int, error) {
	args := i.mock.Called(ctx, q, concurrency)
	err := args.Error(1)
	if err != nil {
		return args.Int(0), err
	}
	return i.client.DeleteWithQuery(ctx, q, concurrency)
}
