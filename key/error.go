package key

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	InvalidPrivateKeyErr = iota
	B58DecodePubKeyErr
	B58DecodeSigErr
	B58ValidateErr
	InvalidDataValidateErr
	SignDataB58Err
	InvalidDataSignErr
	InvalidVerificationKeyErr
	DecodeFromStringErr
	SignError
	JSONError
)

// ErrCodeMessage represents a key-related error.
var ErrCodeMessage = map[int]struct {
	Code    int
	Message string
}{
	InvalidPrivateKeyErr:      {-201, "Private key is invalid"},
	B58DecodePubKeyErr:        {-202, "Base58 decode pub key error"},
	B58DecodeSigErr:           {-203, "Base58 decode signature error"},
	B58ValidateErr:            {-204, "Base58 validate data error"},
	InvalidDataValidateErr:    {-205, "Validated base58 data is invalid"},
	SignDataB58Err:            {-206, "Signing B58 data error"},
	InvalidDataSignErr:        {-207, "Signed data is invalid"},
	InvalidVerificationKeyErr: {-208, "Verification key is invalid"},
	DecodeFromStringErr:       {-209, "Decode key set from string error"},
	SignError:                 {-210, "Can not sign data"},
	JSONError:                 {-211, "JSON Marshal, Unmarshal error"},
}

// Error represents a wrapped error when using the key package.
type Error struct {
	Code    int
	Message string
	err     error
}

// Error returns the beautified string message of an Error.
func (e Error) Error() string {
	return fmt.Sprintf("%d: %s %+v", e.Code, e.Message, e.err)
}

// NewError creates a new Error given a code and an error.
func NewError(key int, err error) *Error {
	return &Error{
		err:     errors.Wrap(err, ErrCodeMessage[key].Message),
		Code:    ErrCodeMessage[key].Code,
		Message: ErrCodeMessage[key].Message,
	}
}
