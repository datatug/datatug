package test

import (
	"fmt"
	"github.com/strongo/validation"
	"testing"
)

type record interface {
	Validate() error
}

func ValidRecord(t *testing.T, name string, r record) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		if err := r.Validate(); err != nil {
			t.Error(fmt.Sprintf("unexpected error of type %T for record of type %T: %+v", err, r, r), err)
		}
	})
}

func InvalidRecord(t *testing.T, name string, r record) {
	t.Run(name, func(t *testing.T) {
		err := r.Validate()
		if err == nil {
			t.Errorf("expected an error but got nil for r: %+v", r)
		}
		if !validation.IsValidationError(err) {
			t.Errorf("returned error is not a validation error: %T: %v; record: %+v", err, err, r)
		}
		if !validation.IsBadRecordError(err) {
			t.Errorf("returned error is not a bad record error: %T: %v; record: %+v", err, err, r)
		}
	})
}
