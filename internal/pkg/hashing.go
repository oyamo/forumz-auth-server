package pkg

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

const (
	Iterations uint32 = 3
	Memory     uint32 = 64 * 1024 // KibiBytes
	Paralelism uint8  = 2
	KeyLength  uint32 = 32
	SaltLength uint32 = 16
)

type Params struct {
	Memory     uint32
	Iterations uint32
	Paralelism uint8
	KeyLength  uint32
	SaltLength uint32
}

// Hash Generation steps
// 1. Generate cryptographycally secure random salt
// 2. Pass the plaintext password, salt and parameters to the argon2.IDKey
// 3. Return a string using the standard encoded hash representation.
// Encoding format
// 		algorithm ID (e.g. argon2id)
//		salt
//		number of iterations (4)
//		memory usage factor (18)
//		parallelism (4)

// Password/Hash comparison
// 1. Decode hash
// 2. Hash the other password and derive key
// 3. Check if the contents are identical

func HashPassword(rawPassword string) (string, error) {
	salt, err := generateRandomBytes(SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(rawPassword), []byte(salt), Iterations, Memory, Paralelism, KeyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, Memory, Iterations, Paralelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func ComparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := DecodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Paralelism, p.KeyLength)

	// Check that the contents of the hashed passwords are identical
	// subtle.ConstantTimeCompare() isto help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func DecodeHash(encodedHash string) (p *Params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return p, nil, nil, errors.New("invalid hash format")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return p, nil, nil, err
	}

	if version != argon2.Version {
		return p, nil, nil, errors.New("invalid hash version")
	}
	p = new(Params)
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Paralelism)
	if err != nil {
		return p, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return p, nil, nil, err
	}

	p.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return p, nil, nil, err
	}

	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
