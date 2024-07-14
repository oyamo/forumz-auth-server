package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"os"
	"software.sslmate.com/src/go-pkcs12"
)

func readFile(path string) ([]byte, error) {
	data, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer data.Close()
	buffer := make([]byte, 1024)

	var b []byte
	for {
		n, err := data.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		b = append(b, buffer[:n]...)
	}
	return b, nil
}

func getPrivateKeyFromP12(path, password string) (*rsa.PrivateKey, error) {
	content, err := readFile(path)
	if err != nil {
		return nil, err
	}

	privateKey, _, err := pkcs12.Decode(content, password)
	if err != nil {
		return nil, err
	}

	return privateKey.(*rsa.PrivateKey), nil
}

func getPublicKeyFromFile(path string) (*rsa.PublicKey, error) {
	content, err := readFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(content)
	if block != nil {
		// PEM format
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("not an RSA public key")
		}
		return rsaKey, nil
	}

	// If not PEM, try parsing as DER format
	key, err := x509.ParsePKIXPublicKey(content)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return rsaKey, nil
}
