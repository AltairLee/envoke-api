package crypto

import (
	. "github.com/zbo14/envoke/common"
	"golang.org/x/crypto/bcrypt"
)

// Interfaces

type PrivateKey interface {
	IsPrivateKey()
	Public() PublicKey
	Sign([]byte) Signature
	String() string
}

type PublicKey interface {
	IsPublicKey()
	Bytes() []byte
	Equals(PublicKey) bool
	FromBytes([]byte) error
	FromString(string) error
	MarshalJSON() ([]byte, error)
	String() string
	UnmarshalJSON([]byte) error
	Verify([]byte, Signature) bool
}

type Signature interface {
	IsSignature()
	Bytes() []byte
	Equals(Signature) bool
	FromBytes([]byte) error
	FromString(string) error
	MarshalJSON() ([]byte, error)
	String() string
	UnmarshalJSON([]byte) error
}

// Generate secret from password using bcrypt

func GenerateSecret(password string) []byte {
	secret, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	Check(err)
	return secret
}
