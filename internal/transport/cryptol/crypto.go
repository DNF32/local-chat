package cryptol

type Service interface {
	Encrypt(data []byte) (string, error)
	Decrypt(encryptedData string) ([]byte, error)
}
