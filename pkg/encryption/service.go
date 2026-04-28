package encryption

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

// Service defines the interface for encryption and decryption.
type Service interface {
	Encrypt(ctx context.Context, plaintext string) (string, error)
	Decrypt(ctx context.Context, ciphertext string) (string, error)
}

type kmsService struct {
	client *kms.Client
	keyID  string
}

var (
	ErrInvalidBundleFormat      = errors.New("invalid bundle format")
	ErrInvalidBundleFormatNonce = errors.New("invalid bundle format: missing nonce")
	ErrDecodeBundle             = errors.New("failed to decode bundle")
	ErrDecryptDataKey           = errors.New("failed to decrypt data key")
	ErrCreateCipher             = errors.New("failed to create AES cipher")
	ErrCreateGCM                = errors.New("failed to create GCM")
	ErrDecryptCiphertext        = errors.New("failed to decrypt ciphertext")
	ErrInvalidCypherText        = errors.New("invalid cypher text")
	ErrGenerateDataKey          = errors.New("failed to generate data key")
	ErrCreateAESCipher          = errors.New("failed to create AES cipher")
	ErrCreateAESGCM             = errors.New("failed to create AES-GCM")
	ErrGenerateNonce            = errors.New("failed to generate nonce")
)

// NewKMSService creates a new encryption service backed by AWS KMS.
func NewKMSService(client *kms.Client, keyID string) Service {
	return &kmsService{
		client: client,
		keyID:  keyID,
	}
}

// Encrypt performs envelope encryption on the plaintext using a KMS-generated data key.
// The returned string is a base64-encoded bundle of [encryptedDataKey][nonce][ciphertext].
func (s *kmsService) Encrypt(ctx context.Context, plaintext string) (string, error) {
	// 1. Generate a data key from KMS
	out, err := s.client.GenerateDataKey(ctx, &kms.GenerateDataKeyInput{
		KeyId:   aws.String(s.keyID),
		KeySpec: types.DataKeySpecAes256,
	})
	if err != nil {
		return "", fmt.Errorf("Encrypt: %w: %w", ErrGenerateDataKey, err)
	}

	// 2. Encrypt the plaintext with the data key using AES-GCM
	block, err := aes.NewCipher(out.Plaintext)
	if err != nil {
		return "", fmt.Errorf("Encrypt: %w: %w", ErrCreateAESCipher, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("Encrypt: %w: %w", ErrCreateAESGCM, err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("Encrypt: %w: %w", ErrGenerateNonce, err)
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// 3. Package encrypted data key + nonce + ciphertext
	// Format: [len(encKey)][encKey][nonce][ciphertext]
	keyLen := len(out.CiphertextBlob)
	bundle := make([]byte, 1+keyLen+len(nonce)+len(ciphertext))
	bundle[0] = byte(keyLen)
	copy(bundle[1:], out.CiphertextBlob)
	copy(bundle[1+keyLen:], nonce)
	copy(bundle[1+keyLen+len(nonce):], ciphertext)

	return base64.StdEncoding.EncodeToString(bundle), nil
}

func (s *kmsService) Decrypt(ctx context.Context, bundleBase64 string) (string, error) {
	bundle, err := base64.StdEncoding.DecodeString(bundleBase64)
	if err != nil {
		return "", fmt.Errorf("Decrypt: %w: %w", ErrDecodeBundle, err)
	}

	if len(bundle) < 1 {
		return "", fmt.Errorf("Decrypt: %w", ErrInvalidBundleFormat)
	}

	keyLen := int(bundle[0])
	if len(bundle) < 1+keyLen {
		return "", fmt.Errorf("Decrypt: key length=%d: %w", keyLen, ErrInvalidBundleFormat)
	}

	encryptedKey := bundle[1 : 1+keyLen]

	out, err := s.client.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob: encryptedKey,
		KeyId:          aws.String(s.keyID),
	})
	if err != nil {
		return "", fmt.Errorf("Decrypt: %w: %w", ErrDecryptDataKey, err)
	}

	block, err := aes.NewCipher(out.Plaintext)
	if err != nil {
		return "", fmt.Errorf("Decrypt: %w: %w", ErrCreateCipher, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("Decrypt: %w: %w", ErrCreateGCM, err)
	}

	nonceSize := gcm.NonceSize()
	if len(bundle) < 1+keyLen+nonceSize {
		return "", fmt.Errorf(
			"Decrypt: bundle_len=%d key_len=%d nonce_size=%d: %w",
			len(bundle), keyLen, nonceSize, ErrInvalidBundleFormatNonce,
		)
	}

	nonce := bundle[1+keyLen : 1+keyLen+nonceSize]
	ciphertext := bundle[1+keyLen+nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("Decrypt: %w: %w", ErrDecryptCiphertext, err)
	}

	return string(plaintext), nil
}

type mockEncryptionService struct{}

func NewMockEncryptionService() Service {
	return &mockEncryptionService{}
}

func (*mockEncryptionService) Encrypt(_ context.Context, plaintext string) (string, error) {
	return "enc:" + plaintext, nil
}

func (*mockEncryptionService) Decrypt(_ context.Context, ciphertext string) (string, error) {
	if len(ciphertext) < 4 || ciphertext[:4] != "enc:" {
		return "", ErrInvalidCypherText
	}
	return ciphertext[4:], nil
}
