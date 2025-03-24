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

	"cloud.google.com/go/firestore"
)

type innerTran struct {
	client *firestore.Client
	tran   *firestore.Transaction
}

func (i *innerTran) Create(ctx context.Context, data any) error {
	rv := reflect.ValueOf(data)
	if p, ok := rv.Interface().(Pathable); ok {
		return i.tran.Create(i.client.Doc(p.Path(ctx)), data)
	}
	return fmt.Errorf("not implement Pathable")
}

func (i *innerTran) Set(ctx context.Context, data any) error {
	rv := reflect.ValueOf(data)
	if p, ok := rv.Interface().(Pathable); ok {
		return i.tran.Set(i.client.Doc(p.Path(ctx)), data)
	}
	return fmt.Errorf("not implement Pathable")
}

func (i *innerTran) Get(ctx context.Context, data any) error {
	rv := reflect.ValueOf(data)
	if p, ok := rv.Interface().(Pathable); ok {
		snapshot, err := i.tran.Get(i.client.Doc(p.Path(ctx)))
		if err != nil {
			return err
		}
		return snapshot.DataTo(data)
	}
	return fmt.Errorf("not implement Pathable")
}

func (i *innerTran) Delete(ctx context.Context, data any) error {
	rv := reflect.ValueOf(data)
	if p, ok := rv.Interface().(Pathable); ok {
		return i.tran.Delete(i.client.Doc(p.Path(ctx)))
	}
	return fmt.Errorf("not implement Pathable")
}
