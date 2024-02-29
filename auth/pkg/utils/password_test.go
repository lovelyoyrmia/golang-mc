package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := RandomString(10)
	hashPassword, err := HashPassword(password)
	require.NoError(t, err)

	err = ComparePassword(hashPassword, password)
	require.NoError(t, err)
}
