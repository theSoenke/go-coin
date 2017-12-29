package coin

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

const (
	version            = byte(0x00)
	addressChecksumLen = 4
	walletFile         = "wallet.dat"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() (*Wallet, error) {
	private, public, err := newKeyPair()
	if err != nil {
		return nil, err
	}
	wallet := Wallet{private, public}

	return &wallet, nil
}

func newKeyPair() (ecdsa.PrivateKey, []byte, error) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return ecdsa.PrivateKey{}, nil, err
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey, nil
}

func (w Wallet) GetAddress() (string, error) {
	pubKeyHash, err := HashPubKey(w.PublicKey)
	if err != nil {
		return "", err
	}

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return string(address), err
}

// ValidateAddress check if address if valid
func ValidateAddress(address string) bool {
	if address == "" {
		return false
	}

	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

func HashPubKey(pubKey []byte) ([]byte, error) {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		return nil, err
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160, nil
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}
