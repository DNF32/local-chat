package transport

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type SerdeAble interface {
	Encode() ([]byte, error)
	Decode(data []byte) error
}

type CryptoService interface {
	Encrypt(data []byte) (string, error)
	Decrypt(encryptedData string) ([]byte, error)
}

type Serde struct {
	cryptoService CryptoService
	logger        *slog.Logger
}

type EncryptedData struct {
	Data string `json:"data"`
}

func (sd *Serde) EncodeEncrypted(obj SerdeAble) ([]byte, error) {
	data, err := obj.Encode()
	if err != nil {
		return nil, err
	}

	encrypted, err := sd.cryptoService.Encrypt(data)
	if err != nil {
		return nil, err
	}

	encryptedBytes, err := json.Marshal(EncryptedData{Data: encrypted})
	if err != nil {
		return nil, fmt.Errorf("Failed to Marshall encrypted data: %w", err)
	}

	return encryptedBytes, nil
}

func (sd *Serde) DecodeEncrypted(data []byte, obj SerdeAble) error {
	var enData EncryptedData
	err := json.Unmarshal(data, &enData)
	if err != nil {
		return fmt.Errorf("Failed to UnMarshall encrypted data: %w", err)
	}

	decrypted, err := sd.cryptoService.Decrypt(enData.Data)
	if err != nil {
		return fmt.Errorf("Failed to decrypt data: %w", err)
	}

	err = obj.Decode(decrypted)
	if err != nil {
		return fmt.Errorf("Failed to decode data: %w", err)
	}

	return nil
}
