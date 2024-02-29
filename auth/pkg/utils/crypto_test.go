package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCrypto(t *testing.T) {
	randomString := RandomString(10)
	encrypted, err := Encrypt(randomString)
	require.NoError(t, err)

	decrypted, err := Decrypt(encrypted)
	require.NoError(t, err)

	require.Equal(t, randomString, decrypted)
}
