package fs

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/adlerhurst/eventstore/v0"
)

var (
	_      eventstore.TestEventstore = (*testStorage)(nil)
	folder fs.FS
	path   string
)

type testStorage struct {
	*Validator
}

// After implements eventstore.TestEventstore
func (*testStorage) After(ctx context.Context, t testing.TB) error {
	// return os.Remove(path)
	return nil
}

// Before implements eventstore.TestEventstore
func (s *testStorage) Before(ctx context.Context, t testing.TB) (err error) {
	return nil
}

var store *testStorage

func TestMain(m *testing.M) {
	var err error
	path, err = os.MkdirTemp(".", "test")
	if err != nil {
		panic(fmt.Sprintf("unable to create temp dir: %v", err))
	}
	folder = os.DirFS(path)
	store = &testStorage{NewValidator(NewFS(path, folder))}
	os.Exit(m.Run())
}
