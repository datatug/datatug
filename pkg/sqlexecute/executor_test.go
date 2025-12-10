package sqlexecute

import "testing"

func TestNewExecutor(t *testing.T) {
	_ = NewExecutor(nil, nil)
}
