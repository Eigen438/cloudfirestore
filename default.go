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

	"cloud.google.com/go/firestore"
)

var defaultInstance CloudFirestore

// Initialize global instance
func Initialize(ctx context.Context) error {
	c, err := New(ctx)
	if err != nil {
		return err
	}
	defaultInstance = c
	return nil
}

// Create Pathable data
func Create(ctx context.Context, data Pathable) error {
	return defaultInstance.Create(ctx, data)
}

// Write/Set Pathable data
func Set(ctx context.Context, data Pathable) error {
	return defaultInstance.Set(ctx, data)
}

// Read/Get Pathable data
func Get(ctx context.Context, data Pathable) error {
	return defaultInstance.Get(ctx, data)
}

// Delete Pathable data
func Delete(ctx context.Context, data Pathable) error {
	return defaultInstance.Delete(ctx, data)
}

// Run transaction
func RunTransaction(ctx context.Context, f func(context.Context, Transaction) error) error {
	return defaultInstance.RunTransaction(ctx, f)
}

// Get collection query
func Collection(collectionName string) firestore.Query {
	return defaultInstance.Collection(collectionName)
}

// Get collection group query
func CollectionGroup(collectionName string) firestore.Query {
	return defaultInstance.CollectionGroup(collectionName)
}

// Sequence query
func Sequence(ctx context.Context, q firestore.Query, f func(context.Context, *firestore.DocumentSnapshot) error) (int, error) {
	return defaultInstance.Sequence(ctx, q, f)
}

// Run query
func Run(ctx context.Context, q firestore.Query, concurrency int, f func(context.Context, *firestore.DocumentSnapshot) error) (int, error) {
	return defaultInstance.Run(ctx, q, concurrency, f)
}

// Delete with Query
func DeleteWithQuery(ctx context.Context, q firestore.Query, concurrency int) (int, error) {
	return defaultInstance.DeleteWithQuery(ctx, q, concurrency)
}

// Sequence query
func TypeSequence[T any](ctx context.Context, q firestore.Query, f func(ctx context.Context, data *T, ref *firestore.DocumentRef) error) (int, error) {
	return defaultInstance.Sequence(ctx, q, func(ctx context.Context, s *firestore.DocumentSnapshot) error {
		data := new(T)
		if err := s.DataTo(data); err != nil {
			return err
		}
		return f(ctx, data, s.Ref)
	})
}

// Run query
func TypedRun[T any](ctx context.Context, q firestore.Query, concurrency int, f func(ctx context.Context, data *T, ref *firestore.DocumentRef) error) (int, error) {
	return defaultInstance.Run(ctx, q, concurrency, func(ctx context.Context, s *firestore.DocumentSnapshot) error {
		data := new(T)
		if err := s.DataTo(data); err != nil {
			return err
		}
		return f(ctx, data, s.Ref)
	})
}

// Return default instance
func Default() CloudFirestore {
	return defaultInstance
}
