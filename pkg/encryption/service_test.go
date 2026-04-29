package encryption

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockEncryptionService(t *testing.T) {
	service := NewMockEncryptionService()
	ctx := context.Background()
	plaintext := "123-45-6789"

	ciphertext, err := service.Encrypt(ctx, plaintext)
	require.NoError(t, err)
	require.NotEqual(t, plaintext, ciphertext)
	require.Equal(t, "enc:"+plaintext, ciphertext)

	decrypted, err := service.Decrypt(ctx, ciphertext)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted)
}

func TestMockEncryptionService_DecryptError(t *testing.T) {
	service := NewMockEncryptionService()
	ctx := context.Background()

	_, err := service.Decrypt(ctx, "invalid")
	require.Error(t, err)
}
