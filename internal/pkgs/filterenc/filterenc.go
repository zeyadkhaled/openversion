// Package filterenc used to convert filter to cursor and cursor to filter.
package filterenc

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/internal/pkgs/errs"
)

type Encer struct {
	key []byte
}

func New(key []byte) Encer {
	k := make([]byte, len(key))
	copy(k, key)
	return Encer{key: k}
}

func (enc Encer) FilterFromCursor(cursor string, filter interface{}) error {
	b, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return errs.E{
			Kind:       errs.KindParameterErr,
			Parameters: []string{"cursor"},
			Wrapped:    err,
		}
	}
	b, err = aeadDecrypt(enc.key, b)
	if err != nil {
		return errs.E{
			Kind:       errs.KindParameterErr,
			Parameters: []string{"cursor"},
			Wrapped:    err,
		}
	}

	err = json.Unmarshal(b, filter)
	if err != nil {
		return errs.E{
			Kind:       errs.KindParameterErr,
			Parameters: []string{"cursor"},
			Wrapped:    err,
		}
	}
	return nil

}

func (enc Encer) CursorFromFilter(filter interface{}) (string, error) {

	b, err := json.Marshal(filter)
	if err != nil {
		return "", fmt.Errorf("failed to marshal filter: %v", err)
	}

	b, err = aeadEncrypt(enc.key, b)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt cursor: %v", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func aeadEncrypt(key, plaintext []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	sealed := aead.Seal(nonce, nonce[:chacha20poly1305.NonceSizeX], plaintext, nil)

	return sealed, nil

}

func aeadDecrypt(key, ciphertext []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < chacha20poly1305.NonceSizeX {
		return nil, errors.New("invalid ciphertext")
	}
	nonce, cipher := ciphertext[:chacha20poly1305.NonceSizeX], ciphertext[chacha20poly1305.NonceSizeX:]
	return aead.Open(nil, nonce, cipher, nil)
}
