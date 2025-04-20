package mocks

import (
	"github.com/stra1g/saver-api/pkg/hashing"
	"github.com/stretchr/testify/mock"
)

type MockHashing struct {
	mock.Mock
}

func NewMockHashing() *MockHashing {
	return &MockHashing{}
}

func (m *MockHashing) HashValue(value string) (string, error) {
	args := m.Called(value)
	return args.String(0), args.Error(1)
}

func (m *MockHashing) CompareHashAndValue(hash, value string) bool {
	args := m.Called(hash, value)
	return args.Bool(0)
}

// Ensure MockHashing implements hashing.Hashing
var _ hashing.Hashing = (*MockHashing)(nil)
