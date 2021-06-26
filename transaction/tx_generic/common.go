package tx_generic

import (
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

func GetTxMintData(tx metadata.Transaction, tokenID *common.Hash) (bool, coin.Coin, *common.Hash, error) {
	outputCoins, err := tx.GetReceiverData()
	if err != nil {
		return false, nil, nil, err
	}
	if len(outputCoins) != 1 {
		return false, nil, nil, errors.New("Error Tx mint has more than one receiver")
	}
	if inputCoins := tx.GetProof().GetInputCoins(); len(inputCoins) > 0 {
		return false, nil, nil, errors.New("Error this is not Tx mint")
	}
	return true, outputCoins[0], tokenID, nil
}

func GetTxBurnData(tx metadata.Transaction) (bool, coin.Coin, *common.Hash, error) {
	outputCoins, err := tx.GetReceiverData()
	if err != nil {
		return false, nil, nil, err
	}
	// remove rule only accept maximum 2 outputs in tx burn
	//if len(outputCoins) > 2 {
	//	utils.Logger.Log.Error("GetAndCheckBurning receiver: More than 2 receivers")
	//	return false, nil, nil, err
	//}
	for _, coin := range outputCoins {
		if wallet.IsPublicKeyBurningAddress(coin.GetPublicKey().ToBytesS()) {
			return true, coin, &common.PRVCoinID, nil
		}
	}
	return false, nil, nil, nil
}

func CalculateSumOutputsWithFee(outputCoins []coin.Coin, fee uint64) *crypto.Point {
	sumOutputsWithFee := new(crypto.Point).Identity()
	for i := 0; i < len(outputCoins); i += 1 {
		sumOutputsWithFee.Add(sumOutputsWithFee, outputCoins[i].GetCommitment())
	}
	feeCommitment := new(crypto.Point).ScalarMult(
		privacy.PedCom.G[privacy.PedersenValueIndex],
		new(crypto.Scalar).FromUint64(fee),
	)
	sumOutputsWithFee.Add(sumOutputsWithFee, feeCommitment)
	return sumOutputsWithFee
}

func ValidateTxParams(params *TxPrivacyInitParams) error {
	if len(params.InputCoins) > 255 {
		return fmt.Errorf("number of inputs (%v) is too large", len(params.InputCoins))
	}
	if len(params.PaymentInfo) > 254 {
		return fmt.Errorf("number of outputs (%v) is too large", len(params.PaymentInfo))
	}
	if params.TokenID == nil {
		// using default PRV
		params.TokenID = &common.Hash{}
		err := params.TokenID.SetBytes(common.PRVCoinID[:])
		if err != nil {
			return fmt.Errorf("cannot setbytes tokenID %v: %v", params.TokenID.String(), err)
		}
	}
	return nil
}

func ParseTokenID(tokenID *common.Hash) (*common.Hash, error) {
	if tokenID == nil {
		result := new(common.Hash)
		err := result.SetBytes(common.PRVCoinID[:])
		if err != nil {
			return nil, fmt.Errorf("cannot parse tokenID %v: %v", tokenID.String(), err)
		}
		return result, nil
	}
	return tokenID, nil
}

func SignNoPrivacy(privKey *key.PrivateKey, hashedMessage []byte) (signatureBytes []byte, sigPubKey []byte, err error) {
	/****** using Schnorr signature *******/
	// sign with sigPrivKey
	// prepare private key for Schnorr
	sk := new(crypto.Scalar).FromBytesS(*privKey)
	r := new(crypto.Scalar).FromUint64(0)
	sigKey := new(privacy.SchnorrPrivateKey)
	sigKey.Set(sk, r)
	signature, err := sigKey.Sign(hashedMessage)
	if err != nil {
		return nil, nil, err
	}

	signatureBytes = signature.Bytes()
	sigPubKey = sigKey.GetPublicKey().GetPublicKey().ToBytesS()
	return signatureBytes, sigPubKey, nil
}
