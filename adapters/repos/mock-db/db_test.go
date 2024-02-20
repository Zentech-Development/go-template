package mockdb_test

import (
	"testing"

	mockAdapters "github.com/Zentech-Development/go-template/adapters/repos/mock-db"
)

func TestNewMockDB(t *testing.T) {
	_ = mockAdapters.NewMockDB()
}
