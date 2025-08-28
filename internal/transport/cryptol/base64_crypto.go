package cryptol

import "encoding/base64"

type Base64Encoder struct{}

func (b *Base64Encoder) Encrypt(data []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(data), nil
}

func (b *Base64Encoder) Decrypt(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
