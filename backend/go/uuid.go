package gofaas

import (
	"github.com/satori/go.uuid"
)

// UUIDGen is a UUID generator that can be mocked for testing
var UUIDGen = func() uuid.UUID {
	v4, _ := uuid.NewV4()
	return v4
}
