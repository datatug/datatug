package firestoreschema

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"github.com/datatug/datatug-cli/pkg/dbschema"
	"google.golang.org/api/iterator"
)

type GetClient func(ctx context.Context) (client *firestore.Client, err error)

func NewProvider(getConnection GetClient) dbschema.Provider {
	return &provider{
		getClient: getConnection,
	}
}

type provider struct {
	getClient GetClient
}

func (p provider) GetCollection(ctx context.Context, path ...string) (collection *dbschema.Collection, err error) {
	var client *firestore.Client
	if client, err = p.getClient(ctx); err != nil {
		return
	}
	iter := client.Collections(ctx)
	for {
		var ref *firestore.CollectionRef
		if ref, err = iter.Next(); err != nil {
			if errors.Is(err, iterator.Done) {
				err = nil
				break
			}
			return
		}
		if ref.ID == path[0] {
			return &dbschema.Collection{
				ID: ref.ID,
			}, nil
		}
	}
	return
}

func (p provider) GetCollections(ctx context.Context, path ...string) (collections []*dbschema.Collection, err error) {
	var client *firestore.Client
	if client, err = p.getClient(ctx); err != nil {
		return
	}
	defer func() {
		_ = client.Close()
	}()
	iter := client.Collections(ctx)
	for {
		var ref *firestore.CollectionRef
		if ref, err = iter.Next(); err != nil {
			if errors.Is(err, iterator.Done) {
				err = nil
				break
			}
			return
		}
		collections = append(collections, &dbschema.Collection{
			ID: ref.ID,
		})
	}
	return
}
