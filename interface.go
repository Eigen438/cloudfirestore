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

type Pathable interface {
	Path(context.Context) string
}

type CloudFirestore interface {
	// Create creates the document with the given data.
	// It returns an error if a document with the same ID already exists.
	Create(context.Context, any) error
	// Delete deletes the document. If the document doesn't exist, it does nothing
	// and returns no error.
	Delete(context.Context, any) error
	// Get retrieves the document.
	Get(context.Context, any) error
	// Set creates or overwrites the document with the given data.
	Set(context.Context, any) error

	// Run transaction
	RunTransaction(context.Context, func(context.Context, Transaction) error) error

	// Get collection query
	Collection(string) firestore.Query
	// Get collection group query
	CollectionGroup(string) firestore.Query

	// Sequence query
	Sequence(context.Context, firestore.Query, func(context.Context, *firestore.DocumentSnapshot) error) (int, error)
	// Run query
	Run(context.Context, firestore.Query, int, func(context.Context, *firestore.DocumentSnapshot) error) (int, error)

	// Delete with Query
	DeleteWithQuery(context.Context, firestore.Query, int) (int, error)
}

type Transaction interface {
	// Create Pathable data in transaction
	Create(context.Context, any) error
	// Write/Set Pathable data in transaction
	Set(context.Context, any) error
	// Read/Get Pathable data in transaction
	Get(context.Context, any) error
	// Delete Pathable data in transaction
	Delete(context.Context, any) error
}
