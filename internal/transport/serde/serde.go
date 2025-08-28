package serde

import (
	"encoding/json"
	"fmt"
	"local-chat/internal/transport/cryptol"
	"log/slog"
)

type SerdeAble interface {
	Encode() ([]byte, error)
	Decode(data []byte) error
}

type EncryptedData struct {
	Data string `json:"data"`
}

type Serde struct {
	cryptoService cryptol.Service
	logger        *slog.Logger
}

func New(cryptoService cryptol.Service, logger *slog.Logger) *Serde {
	if logger == nil {
		logger = slog.Default() // Use default logger if none provided
	}
	return &Serde{
		cryptoService: cryptoService,
		logger:        logger,
	}
}

func (sd *Serde) EncodeEncrypted(obj SerdeAble) ([]byte, error) {
	data, err := obj.Encode()
	if err != nil {
		return nil, &EncodeError{err}
	}

	encrypted, err := sd.cryptoService.Encrypt(data)
	if err != nil {
		return nil, &EncryptError{err}
	}

	encryptedBytes, err := json.Marshal(EncryptedData{Data: encrypted})
	if err != nil {
		return nil, &EncodeError{fmt.Errorf("marshal error: %w", err)}
	}

	return encryptedBytes, nil
}

func (sd *Serde) DecodeEncrypted(data []byte, obj SerdeAble) error {
	var enData EncryptedData
	err := json.Unmarshal(data, &enData)
	if err != nil {
		return &DecodeError{fmt.Errorf("unmarshal error: %w", err)}
	}

	decrypted, err := sd.cryptoService.Decrypt(enData.Data)
	if err != nil {
		return &DecryptError{err}
	}

	err = obj.Decode(decrypted)
	if err != nil {
		return &DecodeError{err}
	}

	return nil
}
