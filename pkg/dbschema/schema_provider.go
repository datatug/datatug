package dbschema

import "context"

type Provider interface {
	GetCollection(ctx context.Context, path ...string) (*Collection, error)
	GetCollections(ctx context.Context, path ...string) ([]*Collection, error)
}

type Collection struct {
	ID string
}
