package cloudfirestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

type innerTran struct {
	client *firestore.Client
	tran   *firestore.Transaction
}

func (i *innerTran) Create(ctx context.Context, data KeyGenerator) error {
	return i.tran.Create(i.client.Doc(data.GenerateKey(ctx)), data)
}

func (i *innerTran) Set(ctx context.Context, data KeyGenerator) error {
	return i.tran.Set(i.client.Doc(data.GenerateKey(ctx)), data)
}

func (i *innerTran) Get(ctx context.Context, data KeyGenerator) error {
	snapshot, err := i.tran.Get(i.client.Doc(data.GenerateKey(ctx)))
	if err != nil {
		return err
	}
	return snapshot.DataTo(data)
}

func (i *innerTran) Delete(ctx context.Context, data KeyGenerator) error {
	return i.tran.Delete(i.client.Doc(data.GenerateKey(ctx)))
}
