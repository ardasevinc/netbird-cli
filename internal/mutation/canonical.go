package mutation

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

func CanonicalJSON(data []byte) ([]byte, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("decode JSON: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("JSON contains more than one value")
		}
		return nil, fmt.Errorf("decode trailing JSON: %w", err)
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode canonical JSON: %w", err)
	}
	return encoded, nil
}

func Digest(parts ...[]byte) (string, error) {
	hash := sha256.New()
	for _, part := range parts {
		canonical, err := CanonicalJSON(part)
		if err != nil {
			return "", err
		}
		_, _ = hash.Write([]byte{0})
		_, _ = hash.Write(canonical)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func Equivalent(left, right []byte) (bool, error) {
	leftCanonical, err := CanonicalJSON(left)
	if err != nil {
		return false, err
	}
	rightCanonical, err := CanonicalJSON(right)
	if err != nil {
		return false, err
	}
	return bytes.Equal(leftCanonical, rightCanonical), nil
}
