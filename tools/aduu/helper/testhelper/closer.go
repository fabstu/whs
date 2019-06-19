package testhelper

import (
	"testing"
)

type Closer interface {
	Close() error
}

func CloseSurely(t *testing.T, closer Closer) {
	if err := closer.Close(); err != nil {
		t.Fatalf("failed to close: %v", err)
	}
}