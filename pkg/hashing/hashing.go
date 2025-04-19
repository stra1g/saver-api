package hashing

import "golang.org/x/crypto/bcrypt"

type Hashing interface {
	HashValue(value string) (string, error)
	CompareHashAndValue(hash, value string) bool
}

type hashing struct{}

func (h *hashing) HashValue(value string) (string, error) {
	hashedValue, err := bcrypt.GenerateFromPassword([]byte(value), 10)

	if err != nil {
		return "", err
	}

	return string(hashedValue), nil
}

func (h *hashing) CompareHashAndValue(hash, value string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
	return err == nil
}

func NewHashing() Hashing {
	return &hashing{}
}
