package metadata

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"strconv"
)

type MetadataBase struct {
	Type int
}

func (mb *MetadataBase) Sign(privateKey *key.PrivateKey, tx Transaction) error {
	return nil
}

type MetadataBaseWithSignature struct {
	MetadataBase
	Sig []byte `json:"Sig,omitempty"`
}

func NewMetadataBaseWithSignature(thisType int) *MetadataBaseWithSignature {
	return &MetadataBaseWithSignature{MetadataBase: MetadataBase{Type: thisType}, Sig: []byte{}}
}

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

func NewMetadataBase(thisType int) *MetadataBase {
	return &MetadataBase{Type: thisType}
}

func (mb *MetadataBase) CalculateSize() uint64 {
	return 0
}

func (mb *MetadataBase) Validate() error {
	return nil
}

func (mb *MetadataBase) Process() error {
	return nil
}

func (mb MetadataBase) GetType() int {
	return mb.Type
}

func (mb MetadataBase) Hash() *common.Hash {
	record := strconv.Itoa(mb.Type)
	data := []byte(record)
	hash := common.HashH(data)
	return &hash
}

func (mb MetadataBase) HashWithoutSig() *common.Hash {
	return mb.Hash()
}
