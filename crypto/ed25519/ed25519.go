package ed25519

import (
	"bytes"
	. "github.com/zballs/go_resonate/util"
	bcrypt "golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ed25519"
)

const (
	PRIVKEY_SIZE   = ed25519.PrivateKeySize
	PUBKEY_SIZE    = ed25519.PublicKeySize
	SIGNATURE_SIZE = ed25519.SignatureSize
)

type PublicKey struct {
	data ed25519.PublicKey
}

type PrivateKey struct {
	data ed25519.PrivateKey
}

type Signature struct {
	data []byte
}

func NewPrivateKey(data ed25519.PrivateKey) (*PrivateKey, error) {
	if size := len(data); size != PRIVKEY_SIZE {
		return nil, Errorf("Expected privkey with size=%d; got size=%d\n", PRIVKEY_SIZE, size)
	}
	return &PrivateKey{data}, nil
}

func NewPublicKey(data ed25519.PublicKey) (*PublicKey, error) {
	if size := len(data); size != PUBKEY_SIZE {
		return nil, Errorf("Expected pubkey with size=%d; got size=%d\n", PUBKEY_SIZE, size)
	}
	return &PublicKey{data}, nil
}

func NewSignature(data []byte) (*Signature, error) {
	if size := len(data); size != SIGNATURE_SIZE {
		return nil, Errorf("Expected signature with size=%d; got size=%d\n", SIGNATURE_SIZE, size)
	}
	return &Signature{data}, nil
}

func GenerateKeypair(password string) (*PrivateKey, *PublicKey) {
	secret := GenerateSecret(password)
	buf := new(bytes.Buffer)
	buf.Write(secret)
	pub_data, priv_data, err := ed25519.GenerateKey(buf)
	Check(err)
	priv, err := NewPrivateKey(priv_data)
	Check(err)
	pub, err := NewPublicKey(pub_data)
	Check(err)
	return priv, pub
}

// Private Key

func (priv *PrivateKey) Sign(message []byte) *Signature {
	data := ed25519.Sign(priv.data, message)
	sig, err := NewSignature(data)
	Check(err)
	return sig
}

func (priv *PrivateKey) Bytes() []byte {
	return priv.data[:]
}

func (priv *PrivateKey) Public() *PublicKey {
	data := priv.data.Public().(ed25519.PublicKey)
	pub, err := NewPublicKey(data)
	Check(err)
	return pub
}

func (priv *PrivateKey) String() string {
	return BytesToB58(priv.Bytes())
}

func (priv *PrivateKey) FromString(b58 string) error {
	data := BytesFromB58(b58)
	if size := len(data); size != PRIVKEY_SIZE {
		return Errorf("Expected privkey with size=%d; got size=%d\n", PRIVKEY_SIZE, size)
	}
	priv.data = make([]byte, PRIVKEY_SIZE)
	copy(priv.data, data)
	return nil
}

func (priv *PrivateKey) MarshalJSON() ([]byte, error) {
	b58 := priv.String()
	data := MustMarshalJSON(b58)
	return data, nil
}

func (priv *PrivateKey) UnmarshalJSON(data []byte) error {
	var b58 string
	if err := UnmarshalJSON(data, &b58); err != nil {
		return err
	}
	return priv.FromString(b58)
}

// Public Key

func (pub *PublicKey) Verify(message []byte, sig *Signature) bool {
	return ed25519.Verify(pub.data, message, sig.data)
}

func (pub *PublicKey) Bytes() []byte {
	return pub.data[:]
}

func (pub *PublicKey) String() string {
	return BytesToB58(pub.Bytes())
}

func (pub *PublicKey) FromString(b58 string) error {
	data := BytesFromB58(b58)
	if size := len(data); size != PUBKEY_SIZE {
		return Errorf("Expected pubkey with size=%d; got size=%d\n", PUBKEY_SIZE, size)
	}
	pub.data = make([]byte, PUBKEY_SIZE)
	copy(pub.data, data)
	return nil
}

func (pub *PublicKey) MarshalJSON() ([]byte, error) {
	b58 := pub.String()
	data := MustMarshalJSON(b58)
	return data, nil
}

func (pub *PublicKey) UnmarshalJSON(data []byte) error {
	var b58 string
	if err := UnmarshalJSON(data, &b58); err != nil {
		return err
	}
	return pub.FromString(b58)
}

// Signature

func (sig *Signature) Bytes() []byte {
	return sig.data[:]
}

func (sig *Signature) String() string {
	return BytesToB58(sig.Bytes())
}

func (sig *Signature) FromString(b58 string) error {
	data := BytesFromB58(b58)
	if size := len(data); size != SIGNATURE_SIZE {
		return Errorf("Expected signature with size=%d; got size=%d\n", SIGNATURE_SIZE, size)
	}
	sig.data = make([]byte, SIGNATURE_SIZE)
	copy(sig.data, data)
	return nil
}

func (sig *Signature) MarshalJSON() ([]byte, error) {
	b58 := sig.String()
	data := MustMarshalJSON(b58)
	return data, nil
}

func (sig *Signature) UnmarshalJSON(data []byte) error {
	var b58 string
	if err := UnmarshalJSON(data, &b58); err != nil {
		return err
	}
	return sig.FromString(b58)
}

// Generate secret from password string

func GenerateSecret(password string) []byte {
	secret, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	Check(err)
	return secret
}
