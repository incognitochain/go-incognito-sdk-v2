package coin

import (
	"fmt"
	"log"
	"testing"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

func newRandomCoinCA(info *PaymentInfo, tokenID *common.Hash) (*CoinV2, error) {
	pk, err := new(crypto.Point).FromBytesS(info.PaymentAddress.Pk)
	if err != nil {
		return nil, err
	}
	pkBytes := pk.ToBytesS()

	c := new(CoinV2).Init()
	c.SetAmount(new(crypto.Scalar).FromUint64(info.Amount))
	c.SetRandomness(crypto.RandomScalar())
	c.SetSharedRandom(crypto.RandomScalar()) // r
	c.SetSharedConcealRandom(crypto.RandomScalar())
	c.SetInfo(info.Message)

	targetShardID := common.GetShardIDFromLastByte(pkBytes[len(pkBytes)-1])
	// Increase index until have the right shardID
	index := uint32(0)
	publicOTA := info.PaymentAddress.GetOTAPublicKey() // For generating one-time-address
	if publicOTA == nil {
		return nil, fmt.Errorf("public OTA from payment address is nil")
	}
	publicSpend := info.PaymentAddress.GetPublicSpend() // General public key
	// publicView := info.PaymentAddress.GetPublicView() // For generating asset tag and concealing output coin

	rK := new(crypto.Point).ScalarMult(publicOTA, c.GetSharedRandom())
	for i := MaxTriesOTA; i > 0; i-- {
		index++

		// Get publickey
		hash := crypto.HashToScalar(append(rK.ToBytesS(), common.Uint32ToBytes(index)...))
		HrKG := new(crypto.Point).ScalarMultBase(hash)
		publicKey := new(crypto.Point).Add(HrKG, publicSpend)
		c.SetPublicKey(publicKey)

		currentShardID, err := c.GetShardID()
		if err != nil {
			return nil, err
		}
		if currentShardID == targetShardID {
			otaSharedRandomPoint := new(crypto.Point).ScalarMultBase(c.GetSharedRandom())
			concealSharedRandomPoint := new(crypto.Point).ScalarMultBase(c.GetSharedConcealRandom())
			c.SetTxRandomDetail(concealSharedRandomPoint, otaSharedRandomPoint, index)

			rAsset := new(crypto.Point).ScalarMult(publicOTA, c.GetSharedRandom())
			blinder, _ := ComputeAssetTagBlinder(rAsset)
			if tokenID == nil {
				return nil, fmt.Errorf("cannot create coin without tokenID")
			}
			assetTag := crypto.HashToPoint(tokenID[:])
			assetTag.Add(assetTag, new(crypto.Point).ScalarMult(crypto.PedCom.G[PedersenRandomnessIndex], blinder))
			c.SetAssetTag(assetTag)
			// fmt.Printf("Shared secret is %s\n", string(rK.MarshalText()))
			// fmt.Printf("Blinder is %s\n", string(blinder.MarshalText()))
			// fmt.Printf("Asset tag is %s\n", string(assetTag.MarshalText()))
			com, err := c.ComputeCommitmentCA()
			if err != nil {
				return nil, fmt.Errorf("cannot compute commitment for confidential asset")
			}
			c.SetCommitment(com)

			return c, nil
		}
	}

	return nil, fmt.Errorf("cannot create OTA after %v attempts", MaxTriesOTA)
}

func makeTokenID(numTokens int) ([]*common.Hash, map[string]*common.Hash) {
	tokenIdList := make([]*common.Hash, 0)
	rawAssetMaps := make(map[string]*common.Hash)

	tokenIdList = append(tokenIdList, &common.PRVCoinID)
	rawAssetMaps[crypto.HashToPoint(common.PRVCoinID[:]).String()] = &common.PRVCoinID

	for i := 0; i < numTokens; i++ {
		tmpTokenID := common.HashH(common.RandBytes(32))
		tokenIdList = append(tokenIdList, &tmpTokenID)
		rawAssetMaps[crypto.HashToPoint(tmpTokenID[:]).String()] = &tmpTokenID
	}

	return tokenIdList, rawAssetMaps
}

func TestCoinV2_GetTokenId(t *testing.T) {
	tokenIDList, rawAssetTags := makeTokenID(1000)

	for i := 0; i < 1000; i++ {
		shardID := byte(common.RandInt() % common.MaxShardNumber)
		keyWallet, err := wallet.GenRandomWalletForShardID(shardID)
		if err != nil {
			panic(err)
		}
		keySet := keyWallet.KeySet

		paymentInfo := InitPaymentInfo(keySet.PaymentAddress, common.RandUint64(), nil)
		tokenId := tokenIDList[common.RandInt()%len(tokenIDList)]

		randomCoin, err := newRandomCoinCA(paymentInfo, tokenId)
		if err != nil {
			panic(err)
		}

		recoveredTokenId, err := randomCoin.GetTokenId(&keySet, rawAssetTags)
		if err != nil {
			panic(err)
		}

		if recoveredTokenId.String() != tokenId.String() {
			panic(fmt.Errorf("require the tokenId to be %v, got %v", tokenId.String(), recoveredTokenId.String()))
		}

		if i%10 == 0 {
			log.Printf("Finished %v\n", i)
		}
	}

}
