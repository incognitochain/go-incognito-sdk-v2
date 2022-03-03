package common

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"strconv"
)

// MetadataBase is the base field for a Metadata.
type MetadataBase struct {
	Type int
}

// Sign does nothing.
func (mb *MetadataBase) Sign(_ *key.PrivateKey, _ Transaction) error {
	return nil
}

// MetadataBaseWithSignature is the base field for a Metadata that requires authenticity (i.e, has a signature).
type MetadataBaseWithSignature struct {
	MetadataBase
	Sig []byte `json:"Sig,omitempty"`
}

// NewMetadataBaseWithSignature creates a new MetadataBaseWithSignature with the given metadata type.
func NewMetadataBaseWithSignature(thisType int) *MetadataBaseWithSignature {
	return &MetadataBaseWithSignature{MetadataBase: MetadataBase{Type: thisType}, Sig: []byte{}}
}

// Sign signs a Metadata using the provided private key.
func (mbs *MetadataBaseWithSignature) Sign(privateKey *key.PrivateKey, tx Transaction) error {
	hashForMd := tx.HashWithoutMetadataSig()
	if hashForMd == nil {
		// the metadata type does not need signing
		return nil
	}
	if len(mbs.Sig) > 0 {
		return fmt.Errorf("cannot overwrite metadata signature")
	}

	/****** using Schnorr signature *******/
	sk := new(crypto.Scalar).FromBytesS(*privateKey)
	r := new(crypto.Scalar).FromUint64(0)
	sigKey := new(privacy.SchnorrPrivateKey)
	sigKey.Set(sk, r)

	// signing
	signature, err := sigKey.Sign(hashForMd[:])
	if err != nil {
		return err
	}

	// convert signature to byte array
	mbs.Sig = signature.Bytes()
	return nil
}

// NewMetadataBase creates a new MetadataBase with the given metadata type.
func NewMetadataBase(thisType int) *MetadataBase {
	return &MetadataBase{Type: thisType}
}

// CalculateSize returns the size of a metadata in bytes.
func (mb *MetadataBase) CalculateSize() uint64 {
	return 0
}

// GetType returns the type of a metadata.
func (mb MetadataBase) GetType() int {
	return mb.Type
}

// Hash calculates the hash of a metadata.
func (mb MetadataBase) Hash() *common.Hash {
	record := strconv.Itoa(mb.Type)
	data := []byte(record)
	hash := common.HashH(data)
	return &hash
}

// HashWithoutSig calculates the hash of a metadata without including its sig.
func (mb MetadataBase) HashWithoutSig() *common.Hash {
	return mb.Hash()
}
