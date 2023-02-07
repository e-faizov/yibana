package encryption

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"hash"
	"io"
	"os"
)

func ReadPubKey(path string) (*rsa.PublicKey, error) {
	var parsedKey interface{}
	pub, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	pubPem, _ := pem.Decode(pub)
	if pubPem == nil {
		return nil, err
	}
	if pubPem.Type != "RSA PUBLIC KEY" {
		return nil, err
	}

	if parsedKey, err = x509.ParsePKCS1PublicKey(pubPem.Bytes); err != nil {
		return nil, err
	}

	var pubKey *rsa.PublicKey
	var ok bool
	if pubKey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, err
	}
	return pubKey, nil
}

func ReadPrivKey(path string) (*rsa.PrivateKey, error) {
	var parsedKey interface{}
	pub, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	pem, _ := pem.Decode(pub)
	if pem == nil {
		return nil, err
	}
	if pem.Type != "RSA PRIVATE KEY" {
		return nil, err
	}

	if parsedKey, err = x509.ParsePKCS1PrivateKey(pem.Bytes); err != nil {
		return nil, err
	}

	var pubKey *rsa.PrivateKey
	var ok bool
	if pubKey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, err
	}
	return pubKey, nil
}

func EncryptOAEP(hash hash.Hash, random io.Reader, public *rsa.PublicKey, msg []byte, label []byte) ([]byte, error) {
	msgLen := len(msg)
	step := public.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, random, public, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func DecryptOAEP(hash hash.Hash, random io.Reader, private *rsa.PrivateKey, msg []byte, label []byte) ([]byte, error) {
	msgLen := len(msg)
	step := private.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, random, private, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}
