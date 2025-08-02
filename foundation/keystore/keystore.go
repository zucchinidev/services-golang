// Package keystore is a simple key store for the application
// It implements the auth.KeyLookup interface.
// It is an in-memory keystore for JWT support.
package keystore

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

type Key struct {
	privatePEM string
	publicPEM  string
}

// KeyStore represents an in-memory store implementation of the auth.KeyLookup interface
// for use with the auth package.
type KeyStore struct {
	store map[string]Key
}

// New constructs an empty KeyStore ready for use.
func New() *KeyStore {
	return &KeyStore{
		store: make(map[string]Key),
	}
}

// PrivateKey returns the private key for the given key identifier.
// If the key is not found, it returns an error.
func (ks *KeyStore) PrivateKey(kid string) (string, error) {
	key, ok := ks.store[kid]
	if !ok {
		return "", fmt.Errorf("key not found: %s", kid)
	}

	return key.privatePEM, nil
}

// PublicKey returns the public key for the given key identifier.
// If the key is not found, it returns an error.
func (ks *KeyStore) PublicKey(kid string) (string, error) {
	key, ok := ks.store[kid]

	if !ok {
		return "", fmt.Errorf("key not found: %s", kid)
	}

	return key.publicPEM, nil
}

// LoadRSAKeys loads a set of RSA PEM files rooted inside of a directory.
// The name of each PEM file will be used as the key identifier.
// Example: /zarf/keys/dc75a316-e862-45ca-a48b-0d67f229d62b.pem
func (ks *KeyStore) LoadRSAKeys(fsys fs.FS) error {
	fn := func(filename string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("unable to read keys: %w", err)
		}

		if dirEntry.IsDir() {
			return nil
		}

		if filepath.Ext(filename) != ".pem" {
			return nil
		}

		file, err := fsys.Open(filename)
		if err != nil {
			return fmt.Errorf("unable to open key file: %w", err)
		}

		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			return fmt.Errorf("unable to stat key file: %w", err)
		}

		if info.Size() == 0 {
			return fmt.Errorf("key file is empty: %s", filename)
		}

		// limit PEM file size to 1MB. This should be reasonable for
		// almost any PEM file and prevents shenanigans like linking
		// the file to /dev/random or /dev/urandom or something like that.
		pem, err := io.ReadAll(io.LimitReader(file, 1024*1024))
		if err != nil {
			return fmt.Errorf("unable to read privatekey file: %w", err)
		}

		privatePEM := string(pem)
		publicPEM, err := toPublicPEM(privatePEM)
		if err != nil {
			return fmt.Errorf("unable to convert to public PEM: %w", err)
		}

		// remove the .pem extension from the filename
		ks.store[strings.TrimSuffix(filename, filepath.Ext(filename))] = Key{
			privatePEM: privatePEM,
			publicPEM:  publicPEM,
		}

		return nil

	}

	if err := fs.WalkDir(fsys, ".", fn); err != nil {
		return fmt.Errorf("unable to walk directory: %w", err)
	}

	return nil
}

func toPublicPEM(privatePEM string) (string, error) {
	block, _ := pem.Decode([]byte(privatePEM))
	if block == nil {
		return "", fmt.Errorf("unable to decode private key")
	}

	var parsedKey any

	parsedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("unable to parse private key: %w", err)
		}
	}

	pk, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("key is not a valid RSA private key")
	}

	asn1Bytes, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	if err != nil {
		return "", fmt.Errorf("unable to marshal public key: %w", err)
	}

	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	return string(pem.EncodeToMemory(&publicBlock)), nil
}
