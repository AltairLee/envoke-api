package ed25519

import (
	"bytes"
	. "github.com/Envoke-org/envoke-api/common"
	"github.com/Envoke-org/envoke-api/crypto/crypto"
	"golang.org/x/crypto/ed25519"
)

const (
	PRIVKEY_SIZE   = ed25519.PrivateKeySize
	PUBKEY_SIZE    = ed25519.PublicKeySize
	SEED_SIZE      = 32
	SIGNATURE_SIZE = ed25519.SignatureSize
)

type PublicKey struct {
	inner ed25519.PublicKey
}

type PrivateKey struct {
	inner ed25519.PrivateKey
}

type Signature struct {
	p []byte
}

func NewPrivateKey(inner ed25519.PrivateKey) *PrivateKey {
	if len(inner) != PRIVKEY_SIZE {
		panic(ErrInvalidSize)
	}
	return &PrivateKey{inner}
}

func NewPublicKey(inner ed25519.PublicKey) *PublicKey {
	if len(inner) != PUBKEY_SIZE {
		panic(ErrInvalidSize)
	}
	return &PublicKey{inner}
}

func NewSignature(inner []byte) *Signature {
	if len(inner) != SIGNATURE_SIZE {
		panic(ErrInvalidSize)
	}
	return &Signature{inner}
}

func GenerateKeypairFromPassword(password string) (*PrivateKey, *PublicKey) {
	secret := crypto.GenerateSecret(password)
	buf := new(bytes.Buffer)
	buf.Write(secret)
	pubInner, privInner, err := ed25519.GenerateKey(buf)
	Check(err)
	priv := NewPrivateKey(privInner)
	pub := NewPublicKey(pubInner)
	return priv, pub
}

func GenerateKeypairFromSeed(seed []byte) (*PrivateKey, *PublicKey) {
	if len(seed) != SEED_SIZE {
		panic(ErrInvalidSize)
	}
	buf := bytes.NewBuffer(seed)
	pubInner, privInner, err := ed25519.GenerateKey(buf)
	Check(err)
	priv := NewPrivateKey(privInner)
	pub := NewPublicKey(pubInner)
	return priv, pub
}

// Private Key

func (_ *PrivateKey) IsPrivateKey() {}

func (priv *PrivateKey) Bytes() []byte {
	if priv == nil {
		return nil
	}
	return priv.inner
}

func (priv *PrivateKey) FromBytes(p []byte) error {
	if len(p) != PRIVKEY_SIZE {
		return ErrInvalidSize
	}
	priv.inner = make([]byte, PRIVKEY_SIZE)
	copy(priv.inner, p)
	return nil
}

func (priv *PrivateKey) FromString(str string) error {
	p := BytesFromB58(str)
	return priv.FromBytes(p)
}

func (priv *PrivateKey) MarshalJSON() ([]byte, error) {
	if priv == nil {
		return nil, nil
	}
	str := priv.String()
	return MustMarshalJSON(str), nil
}

func (priv *PrivateKey) Public() crypto.PublicKey {
	p := priv.inner.Public().(ed25519.PublicKey)
	return NewPublicKey(p)
}

func (priv *PrivateKey) Sign(message []byte) crypto.Signature {
	p := ed25519.Sign(priv.inner, message)
	return NewSignature(p)
}

func (priv *PrivateKey) String() string {
	return BytesToB58(priv.Bytes())
}

// Public Key
func (_ *PublicKey) IsPublicKey() {}

func (pub *PublicKey) Verify(message []byte, sig crypto.Signature) bool {
	return ed25519.Verify(pub.inner, message, sig.Bytes())
}

func (pub *PublicKey) Bytes() []byte {
	if pub == nil {
		return nil
	}
	return pub.inner
}

func (pub *PublicKey) Equals(other crypto.PublicKey) bool {
	return bytes.Equal(pub.Bytes(), other.Bytes())
}

func (pub *PublicKey) FromBytes(p []byte) error {
	if len(p) != PUBKEY_SIZE {
		return ErrInvalidSize
	}
	pub.inner = make([]byte, PUBKEY_SIZE)
	copy(pub.inner, p)
	return nil
}

func (pub *PublicKey) String() string {
	return BytesToB58(pub.Bytes())
}

func (pub *PublicKey) FromString(str string) error {
	p := BytesFromB58(str)
	return pub.FromBytes(p)
}

func (pub *PublicKey) MarshalJSON() ([]byte, error) {
	if pub == nil {
		return nil, nil
	}
	str := pub.String()
	return MustMarshalJSON(str), nil
}

func (pub *PublicKey) UnmarshalJSON(inner []byte) error {
	var str string
	if err := UnmarshalJSON(inner, &str); err != nil {
		return err
	}
	return pub.FromString(str)
}

// Signature

func (_ *Signature) IsSignature() {}

func (sig *Signature) Bytes() []byte {
	if sig == nil {
		return nil
	}
	return sig.p
}

func (sig *Signature) Equals(other crypto.Signature) bool {
	return bytes.Equal(sig.Bytes(), other.Bytes())
}

func (sig *Signature) FromBytes(p []byte) error {
	if len(p) != SIGNATURE_SIZE {
		return ErrInvalidSize
	}
	sig.p = make([]byte, SIGNATURE_SIZE)
	copy(sig.p, p)
	return nil
}

func (sig *Signature) String() string {
	return BytesToB58(sig.Bytes())
}

func (sig *Signature) FromString(str string) error {
	inner := BytesFromB58(str)
	if len(inner) != SIGNATURE_SIZE {
		return ErrInvalidSize
	}
	sig.p = make([]byte, SIGNATURE_SIZE)
	copy(sig.p, inner)
	return nil
}

func (sig *Signature) MarshalJSON() ([]byte, error) {
	if sig == nil {
		return nil, nil
	}
	str := sig.String()
	return MustMarshalJSON(str), nil
}

func (sig *Signature) UnmarshalJSON(inner []byte) error {
	var str string
	if err := UnmarshalJSON(inner, &str); err != nil {
		return err
	}
	return sig.FromString(str)
}
