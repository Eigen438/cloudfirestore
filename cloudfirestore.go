// MIT License
//
// Copyright (c) 2025 Eigen
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

package cloudfirestore

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"cloud.google.com/go/firestore"
	"go.opentelemetry.io/otel"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type inner struct {
	client *firestore.Client
}

func New(ctx context.Context, opts ...option.ClientOption) (CloudFirestore, error) {
	client, err := firestore.NewClient(ctx, firestore.DetectProjectID, opts...)
	if err != nil {
		return nil, err
	}
	return &inner{
		client: client,
	}, nil
}

func NewWithDatabase(ctx context.Context, databaseID string, opts ...option.ClientOption) (CloudFirestore, error) {
	client, err := firestore.NewClientWithDatabase(ctx, firestore.DetectProjectID, databaseID, opts...)
	if err != nil {
		return nil, err
	}
	return &inner{
		client: client,
	}, nil
}

func (i *inner) Create(ctx context.Context, data any) error {
	tracer := otel.GetTracerProvider().Tracer("cloudfirestore")
	ctx, span := tracer.Start(ctx, "Create("+reflect.TypeOf(data).String()+")")
	defer span.End()

	rv := reflect.ValueOf(data)
	if p, ok := rv.Interface().(Pathable); ok {
		_, err := i.client.Doc(p.Path(ctx)).Create(ctx, data)
		return err
	}
	return fmt.Errorf("not implement Pathable")
}

func (i *inner) Delete(ctx context.Context, data any) error {
	tracer := otel.GetTracerProvider().Tracer("cloudfirestore")
	ctx, span := tracer.Start(ctx, "Delete("+reflect.TypeOf(data).String()+")")
	defer span.End()

	rv := reflect.ValueOf(data)
	if p, ok := rv.Interface().(Pathable); ok {
		_, err := i.client.Doc(p.Path(ctx)).Delete(ctx)
		return err
	}
	return fmt.Errorf("not implement Pathable")
}

func (i *inner) Get(ctx context.Context, data any) error {
	tracer := otel.GetTracerProvider().Tracer("cloudfirestore")
	ctx, span := tracer.Start(ctx, "Get("+reflect.TypeOf(data).String()+")")
	defer span.End()

	rv := reflect.ValueOf(data)
	if p, ok := rv.Interface().(Pathable); ok {
		ss, err := i.client.Doc(p.Path(ctx)).Get(ctx)
		if err != nil {
			return err
		}
		return ss.DataTo(data)
	}
	return fmt.Errorf("not implement Pathable")
}

func (i *inner) Set(ctx context.Context, data any) error {
	tracer := otel.GetTracerProvider().Tracer("cloudfirestore")
	ctx, span := tracer.Start(ctx, "Set("+reflect.TypeOf(data).String()+")")
	defer span.End()

	rv := reflect.ValueOf(data)
	if p, ok := rv.Interface().(Pathable); ok {
		_, err := i.client.Doc(p.Path(ctx)).Set(ctx, data)
		return err
	}
	return fmt.Errorf("not implement Pathable")
}

func (i *inner) RunTransaction(ctx context.Context, f func(context.Context, Transaction) error) error {
	tracer := otel.GetTracerProvider().Tracer("cloudfirestore")
	ctx, span := tracer.Start(ctx, "Transaction/All")
	defer span.End()

	err := i.client.RunTransaction(ctx, func(_ctx context.Context, _t *firestore.Transaction) error {
		_ctx, _span := tracer.Start(_ctx, "Transaction/Process")
		defer _span.End()

		t := &innerTran{
			client: i.client,
			tran:   _t,
		}
		return f(_ctx, t)
	})
	return err
}

func (i *inner) Collection(collectionName string) firestore.Query {
	return i.client.Collection(collectionName).Query
}

func (i *inner) CollectionGroup(collectionName string) firestore.Query {
	return i.client.CollectionGroup(collectionName).Query
}

func (i *inner) Sequence(ctx context.Context, q firestore.Query, f func(ctx context.Context, snapshot *firestore.DocumentSnapshot) error) (int, error) {
	tracer := otel.GetTracerProvider().Tracer("cloudfirestore")
	iter := q.Documents(ctx)
	defer iter.Stop()

	num := 0
	for {
		s, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return num, err
		}

		if err := func(_ctx context.Context, _s *firestore.DocumentSnapshot) error {
			_ctx, span := tracer.Start(_ctx, "Sequence/Process:"+_s.Ref.ID)
			defer span.End()
			return f(_ctx, _s)
		}(ctx, s); err != nil {
			return num, err
		}
		num++

	}
	return num, nil
}

func (i *inner) Run(ctx context.Context, q firestore.Query, concurrency int, f func(ctx context.Context, snapshot *firestore.DocumentSnapshot) error) (int, error) {
	tracer := otel.GetTracerProvider().Tracer("cloudfirestore")
	iter := q.Documents(ctx)
	defer iter.Stop()

	var wg sync.WaitGroup
	var errRet error
	num := 0
	ch := make(chan struct{}, concurrency)
	for {
		s, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			errRet = err
			break
		}

		num++
		wg.Add(1)
		ch <- struct{}{}
		go func(_ctx context.Context, _s *firestore.DocumentSnapshot, _ch chan struct{}) {
			defer wg.Done()

			_ctx, span := tracer.Start(_ctx, "Run/Process:"+_s.Ref.ID)
			defer span.End()
			err = f(_ctx, _s)
			if err != nil {
				errRet = err
			}
			<-_ch
		}(ctx, s, ch)
	}
	wg.Wait()
	return num, errRet
}

func (i *inner) DeleteWithQuery(ctx context.Context, q firestore.Query, concurrency int) (int, error) {
	tracer := otel.GetTracerProvider().Tracer("cloudfirestore")
	ctx, span := tracer.Start(ctx, "Transaction/All")
	defer span.End()

	bw := i.client.BulkWriter(ctx)
	num, err := i.Sequence(ctx, q, func(_ context.Context, snapshot *firestore.DocumentSnapshot) error {
		_, err := bw.Delete(snapshot.Ref)
		return err
	})
	bw.End()
	return num, err
}
