package serde

import "fmt"

// Base interface for serde errors
type SerdeError interface {
	error
	Op() string
	Unwrap() error
}

// EncodeError
type EncodeError struct {
	Err error
}

func (e *EncodeError) Error() string { return fmt.Sprintf("encode error: %v", e.Err) }
func (e *EncodeError) Op() string    { return "encode" }
func (e *EncodeError) Unwrap() error { return e.Err }

// DecodeError
type DecodeError struct {
	Err error
}

func (e *DecodeError) Error() string { return fmt.Sprintf("decode error: %v", e.Err) }
func (e *DecodeError) Op() string    { return "decode" }
func (e *DecodeError) Unwrap() error { return e.Err }

// EncryptError
type EncryptError struct {
	Err error
}

func (e *EncryptError) Error() string { return fmt.Sprintf("encrypt error: %v", e.Err) }
func (e *EncryptError) Op() string    { return "encrypt" }
func (e *EncryptError) Unwrap() error { return e.Err }

// DecryptError
type DecryptError struct {
	Err error
}

func (e *DecryptError) Error() string { return fmt.Sprintf("decrypt error: %v", e.Err) }
func (e *DecryptError) Op() string    { return "decrypt" }
func (e *DecryptError) Unwrap() error { return e.Err }
