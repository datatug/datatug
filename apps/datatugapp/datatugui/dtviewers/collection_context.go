package dtviewers

import "github.com/dal-go/dalgo/dal"

type CollectionContext struct {
	DbContext
	CollectionRef dal.CollectionRef
}
