package token

import (
	"fmt"
	"time"

	"github.com/Foedie/foedie-server-v2/auth/pkg/config"
	"github.com/Foedie/foedie-server-v2/auth/pkg/constants"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto      *paseto.V2
	symetricKey []byte
}

func NewPasetoMaker(conf config.Config) (Maker, error) {
	if len(conf.SecretKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d chars", chacha20poly1305.KeySize)
	}

	return &PasetoMaker{
		paseto:      paseto.NewV2(),
		symetricKey: []byte(conf.SecretKey),
	}, nil
}

// GenerateToken implements Maker.
func (maker *PasetoMaker) GenerateToken(uid string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(uid, duration)
	if err != nil {
		return "", payload, err
	}
	token, err := maker.paseto.Encrypt(maker.symetricKey, payload, nil)
	return token, payload, err
}

// VerifyToken implements Maker.
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symetricKey, payload, nil)
	if err != nil {
		return nil, constants.ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
