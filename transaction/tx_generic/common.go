package tx_generic

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// GetTxMintData returns the minting data of a transaction w.r.t a given tokenID.
func GetTxMintData(tx metadata.Transaction, tokenID *common.Hash) (bool, coin.Coin, *common.Hash, error) {
	outputCoins, err := tx.GetReceiverData()
	if err != nil {
		return false, nil, nil, err
	}
	if len(outputCoins) != 1 {
		return false, nil, nil, fmt.Errorf("tx mint has more than one receiver")
	}
	if inputCoins := tx.GetProof().GetInputCoins(); len(inputCoins) > 0 {
		return false, nil, nil, fmt.Errorf("this is not a tx mint")
	}
	return true, outputCoins[0], tokenID, nil
}

// GetTxBurnData returns the burning data of a transaction.
func GetTxBurnData(tx metadata.Transaction) (bool, coin.Coin, *common.Hash, error) {
	outputCoins, err := tx.GetReceiverData()
	if err != nil {
		return false, nil, nil, err
	}

	for _, c := range outputCoins {
		if wallet.IsPublicKeyBurningAddress(c.GetPublicKey().ToBytesS()) {
			return true, c, &common.PRVCoinID, nil
		}
	}
	return false, nil, nil, nil
}

// CalculateSumOutputsWithFee returns a sum of given output coins' commitments plus the commitment of the given fee.
func CalculateSumOutputsWithFee(outputCoins []coin.Coin, fee uint64) *crypto.Point {
	sumOutputsWithFee := new(crypto.Point).Identity()
	for i := 0; i < len(outputCoins); i += 1 {
		sumOutputsWithFee.Add(sumOutputsWithFee, outputCoins[i].GetCommitment())
	}
	feeCommitment := new(crypto.Point).ScalarMult(
		crypto.PedCom.G[crypto.PedersenValueIndex],
		new(crypto.Scalar).FromUint64(fee),
	)
	sumOutputsWithFee.Add(sumOutputsWithFee, feeCommitment)
	return sumOutputsWithFee
}

// ValidateTxParams checks sanity of a TxPrivacyInitParams.
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

// SignNoPrivacy signs a message in a non-private manner using the Schnorr signature scheme.
func SignNoPrivacy(privateKey *key.PrivateKey, hashedMessage []byte) (signatureBytes []byte, sigPubKey []byte, err error) {
	sk := new(crypto.Scalar).FromBytesS(*privateKey)
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
