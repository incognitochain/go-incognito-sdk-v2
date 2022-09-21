package coin

import (
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

const (
	MaxTriesOTA    int = 50000
	AssetTagString     = "assettag"
)

// ComputeCommitmentCA computes the commitment of a CoinV2 with its asset token.
//
// This function is used to provide confidential assets of a token transaction V2.
func (c *CoinV2) ComputeCommitmentCA() (*crypto.Point, error) {
	if c == nil || c.GetRandomness() == nil || c.GetAmount() == nil {
		return nil, fmt.Errorf("missing arguments for committing")
	}
	randomnessBasePoint := crypto.PedCom.G[crypto.PedersenRandomnessIndex]
	commitment := new(crypto.Point).ScalarMult(c.GetAssetTag(), c.GetAmount())
	commitment.Add(commitment, new(crypto.Point).ScalarMult(randomnessBasePoint, c.GetRandomness()))
	return commitment, nil
}

// ComputeCommitmentCA computes a commitment based on the given assetTag, randomness r and value v.
func ComputeCommitmentCA(assetTag *crypto.Point, r, v *crypto.Scalar) (*crypto.Point, error) {
	if assetTag == nil || r == nil || v == nil {
		return nil, fmt.Errorf("missing arguments for committing to CA coin")
	}
	randomnessBasePoint := crypto.PedCom.G[crypto.PedersenRandomnessIndex]
	commitment := new(crypto.Point).ScalarMult(assetTag, v)
	commitment.Add(commitment, new(crypto.Point).ScalarMult(randomnessBasePoint, r))
	return commitment, nil
}

// ComputeAssetTagBlinder returns the asset tag blinder from a shared secret.
func ComputeAssetTagBlinder(sharedSecret *crypto.Point) (*crypto.Scalar, error) {
	if sharedSecret == nil {
		return nil, fmt.Errorf("missing arguments for asset tag blinder")
	}
	blinder := crypto.HashToScalar(append(sharedSecret.ToBytesS(), []byte(AssetTagString)...))
	return blinder, nil
}

// RecomputeSharedSecret uses the privateKey to re-compute the shared OTA secret of a CoinV2.
func (c *CoinV2) RecomputeSharedSecret(privateKey []byte) (*crypto.Point, error) {
	// sk := new(crypto.Scalar).FromBytesS(privateKey)
	var privateOTA []byte = key.GeneratePrivateOTAKey(privateKey)[:]
	sk := new(crypto.Scalar).FromBytesS(privateOTA)
	// this is g^SharedRandom, previously created by sender of the coin
	sharedOTARandomPoint, err := c.GetTxRandom().GetTxOTARandomPoint()
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve tx random detail")
	}
	sharedSecret := new(crypto.Point).ScalarMult(sharedOTARandomPoint, sk)
	return sharedSecret, nil
}

// ValidateAssetTag checks if the asset tag of a CoinV2 is valid for the given shared secret and tokenID.
func (c *CoinV2) ValidateAssetTag(sharedSecret *crypto.Point, tokenID *common.Hash) (bool, error) {
	if c.GetAssetTag() == nil {
		if tokenID == nil || *tokenID == common.PRVCoinID {
			// a valid PRV coin
			return true, nil
		}
		return false, fmt.Errorf("CA coin must have asset tag")
	}
	if tokenID == nil || *tokenID == common.PRVCoinID {
		// invalid
		return false, fmt.Errorf("PRV coin cannot have asset tag")
	}
	recomputedAssetTag := crypto.HashToPoint(tokenID[:])
	if crypto.IsPointEqual(recomputedAssetTag, c.GetAssetTag()) {
		return true, nil
	}

	blinder, err := ComputeAssetTagBlinder(sharedSecret)
	if err != nil {
		return false, err
	}

	recomputedAssetTag.Add(recomputedAssetTag, new(crypto.Point).ScalarMult(crypto.PedCom.G[PedersenRandomnessIndex], blinder))
	if crypto.IsPointEqual(recomputedAssetTag, c.GetAssetTag()) {
		return true, nil
	}
	return false, nil
}

// SetPlainTokenID sets the given tokenID as the tokenID of a CoinV2 (in raw value, not blinded).
func (c *CoinV2) SetPlainTokenID(tokenID *common.Hash) error {
	assetTag := crypto.HashToPoint(tokenID[:])
	c.SetAssetTag(assetTag)
	com, err := c.ComputeCommitmentCA()
	if err != nil {
		return err
	}
	c.SetCommitment(com)
	return nil
}

// NewCoinCAFromOTAReceiver creates a new CoinV2 with asset tag from a given OTA receiver, amount, info and tokenID.
func NewCoinCAFromOTAReceiver(otaReceiver OTAReceiver, amount uint64, info []byte, tokenID *common.Hash) (*CoinV2, *crypto.Point, error) {
	c := new(CoinV2).Init()
	c.SetPublicKey(&otaReceiver.PublicKey)
	c.SetAmount(new(crypto.Scalar).FromUint64(amount))
	c.SetRandomness(crypto.RandomScalar())
	c.SetTxRandom(&otaReceiver.TxRandom)

	rAsset := &otaReceiver.SharedSecrets[0]
	blinder, _ := ComputeAssetTagBlinder(rAsset)
	assetTag := crypto.HashToPoint(tokenID[:])
	assetTag.Add(assetTag, new(crypto.Point).ScalarMult(crypto.PedCom.G[PedersenRandomnessIndex], blinder))
	c.SetAssetTag(assetTag)

	com, err := c.ComputeCommitmentCA()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot compute commitment for confidential asset")
	}
	c.SetCommitment(com)

	c.SetSharedRandom(nil)
	c.SetSharedConcealRandom(nil)
	c.SetInfo(info)

	return c, rAsset, nil
}

// NewCoinCA creates a new CoinV2 for the paymentInfo with asset tag for the given tokenID.
// It is used in the case of confidential assets only
func NewCoinCA(info *key.PaymentInfo, tokenID *common.Hash) (*CoinV2, *crypto.Point, error) {
	if info.OTAReceiver != "" {
		otaReceiver := new(OTAReceiver)
		err := otaReceiver.FromString(info.OTAReceiver)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid otaReceiver %v: %v", info.OTAReceiver, err)
		}

		c, ss, err := NewCoinCAFromOTAReceiver(*otaReceiver, info.Amount, info.Message, tokenID)
		if err != nil {
			return nil, nil, err
		}

		return c, ss, nil
	}

	receiverPublicKey, err := new(crypto.Point).FromBytesS(info.PaymentAddress.Pk)
	if err != nil {
		errStr := fmt.Sprintf("Cannot parse outputCoinV2 from PaymentInfo when parseByte PublicKey, error %v ", err)
		return nil, nil, fmt.Errorf(errStr)
	}
	receiverPublicKeyBytes := receiverPublicKey.ToBytesS()
	targetShardID := common.GetShardIDFromLastByte(receiverPublicKeyBytes[len(receiverPublicKeyBytes)-1])

	c := new(CoinV2).Init()

	// Amount, Randomness, SharedRandom are transparent until we call concealData
	c.SetAmount(new(crypto.Scalar).FromUint64(info.Amount))
	c.SetRandomness(crypto.RandomScalar())
	c.SetSharedRandom(crypto.RandomScalar()) // r
	c.SetSharedConcealRandom(crypto.RandomScalar())
	c.SetInfo(info.Message)

	// If this is going to burning address then don't need to create ota
	if wallet.IsPublicKeyBurningAddress(info.PaymentAddress.Pk) {
		publicKey, err := new(crypto.Point).FromBytesS(info.PaymentAddress.Pk)
		if err != nil {
			panic("something is wrong with info.paymentAddress.pk, burning address should be a valid point")
		}
		c.SetPublicKey(publicKey)
		err = c.SetPlainTokenID(tokenID)
		if err != nil {
			return nil, nil, err
		}
		return c, nil, nil
	}

	// Increase index until have the right shardID
	index := uint32(0)
	publicOTA := info.PaymentAddress.GetOTAPublicKey() // For generating one-time-address
	if publicOTA == nil {
		return nil, nil, fmt.Errorf("public OTA from payment address is nil")
	}
	publicSpend := info.PaymentAddress.GetPublicSpend() // public key

	rK := new(crypto.Point).ScalarMult(publicOTA, c.GetSharedRandom())
	for i := MaxTriesOTA; i > 0; i-- {
		index++

		hash := crypto.HashToScalar(append(rK.ToBytesS(), common.Uint32ToBytes(index)...))
		HrKG := new(crypto.Point).ScalarMultBase(hash)
		publicKey := new(crypto.Point).Add(HrKG, publicSpend)
		c.SetPublicKey(publicKey)

		senderShardID, receivingShardID, coinPrivacyType, _ := DeriveShardInfoFromCoin(publicKey.ToBytesS())
		if receivingShardID == int(targetShardID) && senderShardID == p.SenderShardID && coinPrivacyType == p.CoinPrivacyType {
			otaSharedRandomPoint := new(crypto.Point).ScalarMultBase(c.GetSharedRandom())
			concealSharedRandomPoint := new(crypto.Point).ScalarMultBase(c.GetSharedConcealRandom())
			c.SetTxRandomDetail(concealSharedRandomPoint, otaSharedRandomPoint, index)

			rAsset := new(crypto.Point).ScalarMult(publicOTA, c.GetSharedRandom())
			blinder, _ := ComputeAssetTagBlinder(rAsset)
			if tokenID == nil {
				return nil, nil, fmt.Errorf("cannot create coin without tokenID")
			}
			assetTag := crypto.HashToPoint(tokenID[:])
			assetTag.Add(assetTag, new(crypto.Point).ScalarMult(crypto.PedCom.G[PedersenRandomnessIndex], blinder))
			c.SetAssetTag(assetTag)
			com, err := c.ComputeCommitmentCA()
			if err != nil {
				return nil, nil, fmt.Errorf("cannot compute commitment for confidential asset")
			}
			c.SetCommitment(com)

			return c, rAsset, nil
		}
	}

	// MaxAttempts could be exceeded if the OS's RNG or the stateDB is corrupted
	return nil, nil, fmt.Errorf("cannot create OTA after %d attempts", MaxTriesOTA)
}
