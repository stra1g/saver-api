package hashing_test

import (
	"testing"

	"github.com/stra1g/saver-api/pkg/hashing"
	"github.com/stretchr/testify/assert"
)

func TestHashValue(t *testing.T) {
	// Arrange
	h := hashing.NewHashing()
	plainText := "password123"

	// Act
	hashedValue, err := h.HashValue(plainText)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedValue)
	assert.NotEqual(t, plainText, hashedValue)
}

func TestCompareHashAndValue(t *testing.T) {
	// Arrange
	h := hashing.NewHashing()
	plainText := "securepassword"

	// First get a hash
	hashedValue, err := h.HashValue(plainText)
	assert.NoError(t, err)

	// Act & Assert - correct password
	result := h.CompareHashAndValue(hashedValue, plainText)
	assert.True(t, result)

	// Act & Assert - incorrect password
	wrongResult := h.CompareHashAndValue(hashedValue, "wrongpassword")
	assert.False(t, wrongResult)
}

func TestCompareHashAndValue_WithEmptyValues(t *testing.T) {
	// Arrange
	h := hashing.NewHashing()

	// Act & Assert - empty hash
	result1 := h.CompareHashAndValue("", "password")
	assert.False(t, result1)

	// Act & Assert - empty value
	hashedValue, _ := h.HashValue("password")
	result2 := h.CompareHashAndValue(hashedValue, "")
	assert.False(t, result2)

	// Act & Assert - both empty
	result3 := h.CompareHashAndValue("", "")
	assert.False(t, result3)
}

func TestMultipleHashingOfSameValue(t *testing.T) {
	// Arrange
	h := hashing.NewHashing()
	plainText := "testpassword"

	// Act
	hash1, err1 := h.HashValue(plainText)
	hash2, err2 := h.HashValue(plainText)

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)

	// The hashes should be different due to salt
	assert.NotEqual(t, hash1, hash2)

	// But both should validate against the original value
	assert.True(t, h.CompareHashAndValue(hash1, plainText))
	assert.True(t, h.CompareHashAndValue(hash2, plainText))
}
