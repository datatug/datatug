package schemers

import (
	"context"

	"github.com/dal-go/dalgo/dal"
)

type Provider interface {
	GetCollection(ctx context.Context, collectionRef *dal.CollectionRef) (*Collection, error)
	GetCollections(ctx context.Context, parentKey *dal.Key) ([]*Collection, error)
}

type Collection struct {
	ID string
}
