package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmailValid(t *testing.T) {
	email := GenerateRandomEmail(10)
	valid_email := IsEmailValid(email)
	require.Equal(t, true, valid_email)
}
