package jsonresult

import (
	"errors"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"math/big"
	"strconv"
)

// ICoinInfo describes all methods of an RPC output coin.
type ICoinInfo interface {
	GetVersion() uint8
	GetCommitment() *crypto.Point
	GetInfo() []byte
	GetPublicKey() *crypto.Point
	GetValue() uint64
	GetKeyImage() *crypto.Point
	GetRandomness() *crypto.Scalar
	GetShardID() (uint8, error)
	GetSNDerivator() *crypto.Scalar
	GetCoinDetailEncrypted() []byte
	IsEncrypted() bool
	GetTxRandom() *coin.TxRandom
	GetSharedRandom() *crypto.Scalar
	GetSharedConcealRandom() *crypto.Scalar
	GetAssetTag() *crypto.Point
}

// ListOutputCoins is a list of output coins returned by an RPC response.
type ListOutputCoins struct {
	FromHeight uint64               `json:"FromHeight"`
	ToHeight   uint64               `json:"ToHeight"`
	Outputs    map[string][]OutCoin `json:"Outputs"`
}

// OutCoin is a struct to parse raw-data returned by an RPC response into an output coin.
type OutCoin struct {
	Version              string `json:"Version"`
	Index                string `json:"Index"`
	PublicKey            string `json:"PublicKey"`
	Commitment           string `json:"Commitment"`
	CoinCommitment       string `json:"CoinCommitment"`
	SNDerivator          string `json:"SNDerivator"`
	KeyImage             string `json:"KeyImage"`
	Randomness           string `json:"Randomness"`
	Value                string `json:"Value"`
	Info                 string `json:"Info"`
	SharedRandom         string `json:"SharedRandom"`
	SharedConcealRandom  string `json:"SharedConcealRandom"`
	TxRandom             string `json:"TxRandom"`
	CoinDetailsEncrypted string `json:"CoinDetailsEncrypted"`
	AssetTag             string `json:"AssetTag"`
}

// Conceal removes all fields of an OutCoin leaving only the version, commitment and public key.
// It is usually used to enhance privacy before being sent to the remote server.
func (outCoin *OutCoin) Conceal() {
	outCoin.Index = ""
	outCoin.SNDerivator = ""
	outCoin.KeyImage = ""
	outCoin.Randomness = ""
	outCoin.Value = ""
	outCoin.Info = ""
	outCoin.SharedRandom = ""
	outCoin.SharedConcealRandom = ""
	outCoin.TxRandom = ""
	outCoin.CoinDetailsEncrypted = ""
	outCoin.AssetTag = ""
}

// NewOutCoin creates a new OutCoin from the given ICoinInfo.
func NewOutCoin(outCoin ICoinInfo) OutCoin {
	keyImage := ""
	if outCoin.GetKeyImage() != nil && !outCoin.GetKeyImage().IsIdentity() {
		keyImage = base58.Base58Check{}.Encode(outCoin.GetKeyImage().ToBytesS(), common.ZeroByte)
	}

	publicKey := ""
	if outCoin.GetPublicKey() != nil {
		publicKey = base58.Base58Check{}.Encode(outCoin.GetPublicKey().ToBytesS(), common.ZeroByte)
	}

	commitment := ""
	if outCoin.GetCommitment() != nil {
		commitment = base58.Base58Check{}.Encode(outCoin.GetCommitment().ToBytesS(), common.ZeroByte)
	}

	snd := ""
	if outCoin.GetSNDerivator() != nil {
		snd = base58.Base58Check{}.Encode(outCoin.GetSNDerivator().ToBytesS(), common.ZeroByte)
	}

	randomness := ""
	if outCoin.GetRandomness() != nil {
		randomness = base58.Base58Check{}.Encode(outCoin.GetRandomness().ToBytesS(), common.ZeroByte)
	}

	info := ""
	if len(outCoin.GetInfo()) != 0 {
		info = EncodeBase58Check(outCoin.GetInfo())
	} else {
		info = "13PMpZ4"
	}

	result := OutCoin{
		Version:        strconv.FormatUint(uint64(outCoin.GetVersion()), 10),
		PublicKey:      publicKey,
		Value:          strconv.FormatUint(outCoin.GetValue(), 10),
		Info:           info,
		Commitment:     commitment,
		CoinCommitment: commitment,
		SNDerivator:    snd,
		KeyImage:       keyImage,
		Randomness:     randomness,
	}

	if outCoin.GetCoinDetailEncrypted() != nil {
		result.CoinDetailsEncrypted = base58.Base58Check{}.Encode(outCoin.GetCoinDetailEncrypted(), common.ZeroByte)
	}

	if outCoin.GetSharedRandom() != nil {
		result.SharedRandom = base58.Base58Check{}.Encode(outCoin.GetSharedRandom().ToBytesS(), common.ZeroByte)
	}
	if outCoin.GetSharedConcealRandom() != nil {
		result.SharedRandom = base58.Base58Check{}.Encode(outCoin.GetSharedConcealRandom().ToBytesS(), common.ZeroByte)
	}
	if outCoin.GetTxRandom() != nil {
		result.TxRandom = base58.Base58Check{}.Encode(outCoin.GetTxRandom().Bytes(), common.ZeroByte)
	}
	if outCoin.GetAssetTag() != nil {
		result.AssetTag = base58.Base58Check{}.Encode(outCoin.GetAssetTag().ToBytesS(), common.ZeroByte)
	}

	return result
}

// NewCoinFromJsonOutCoin returns an ICoinInfo, and an index from an OutCoin.
func NewCoinFromJsonOutCoin(jsonOutCoin OutCoin) (ICoinInfo, *big.Int, error) {
	var keyImage, pubkey, cm *crypto.Point
	var snd, randomness *crypto.Scalar
	var info []byte
	var err error
	var idx *big.Int
	var sharedRandom, sharedConcealRandom *crypto.Scalar
	var txRandom *coin.TxRandom
	var coinDetailEncrypted *privacy.HybridCipherText
	var assetTag *crypto.Point

	value, ok := math.ParseUint64(jsonOutCoin.Value)
	if !ok {
		return nil, nil, errors.New("Cannot parse value")
	}

	if len(jsonOutCoin.KeyImage) == 0 {
		keyImage = nil
	} else {
		keyImageInBytes, _, err := base58.Base58Check{}.Decode(jsonOutCoin.KeyImage)
		if err != nil {
			return nil, nil, err
		}
		keyImage, err = new(crypto.Point).FromBytesS(keyImageInBytes)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(jsonOutCoin.Commitment) == 0 {
		if len(jsonOutCoin.CoinCommitment) == 0 {
			cm = nil
		} else {
			cmInbytes, _, err := base58.Base58Check{}.Decode(jsonOutCoin.CoinCommitment)
			if err != nil {
				return nil, nil, err
			}
			cm, err = new(crypto.Point).FromBytesS(cmInbytes)
			if err != nil {
				return nil, nil, err
			}
		}
	} else {
		cmInbytes, _, err := base58.Base58Check{}.Decode(jsonOutCoin.Commitment)
		if err != nil {
			return nil, nil, err
		}
		cm, err = new(crypto.Point).FromBytesS(cmInbytes)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(jsonOutCoin.PublicKey) == 0 {
		pubkey = nil
	} else {
		pubkeyInBytes, _, err := base58.Base58Check{}.Decode(jsonOutCoin.PublicKey)
		if err != nil {
			return nil, nil, err
		}
		pubkey, err = new(crypto.Point).FromBytesS(pubkeyInBytes)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(jsonOutCoin.Randomness) == 0 {
		randomness = nil
	} else {
		randomnessInBytes, _, err := base58.Base58Check{}.Decode(jsonOutCoin.Randomness)
		if err != nil {
			return nil, nil, err
		}
		randomness = new(crypto.Scalar).FromBytesS(randomnessInBytes)
	}

	if len(jsonOutCoin.SNDerivator) == 0 {
		snd = nil
	} else {
		sndInBytes, _, err := base58.Base58Check{}.Decode(jsonOutCoin.SNDerivator)
		if err != nil {
			return nil, nil, err
		}
		snd = new(crypto.Scalar).FromBytesS(sndInBytes)
	}

	if len(jsonOutCoin.Info) == 0 {
		info = []byte{}
	} else {
		info, _, err = base58.Base58Check{}.Decode(jsonOutCoin.Info)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(jsonOutCoin.SharedRandom) == 0 {
		sharedRandom = nil
	} else {
		sharedRandomInBytes, _, err := base58.Base58Check{}.Decode(jsonOutCoin.SharedRandom)
		if err != nil {
			return nil, nil, err
		}
		sharedRandom = new(crypto.Scalar).FromBytesS(sharedRandomInBytes)
	}

	if len(jsonOutCoin.SharedConcealRandom) == 0 {
		sharedRandom = nil
	} else {
		sharedConcealRandomInBytes, _, err := base58.Base58Check{}.Decode(jsonOutCoin.SharedConcealRandom)
		if err != nil {
			return nil, nil, err
		}
		sharedConcealRandom = new(crypto.Scalar).FromBytesS(sharedConcealRandomInBytes)
	}

	if len(jsonOutCoin.TxRandom) == 0 {
		sharedRandom = nil
	} else {
		txRandomInBytes, _, err := base58.Base58Check{}.Decode(jsonOutCoin.TxRandom)
		if err != nil {
			return nil, nil, err
		}
		txRandom = new(coin.TxRandom)
		err = txRandom.SetBytes(txRandomInBytes)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(jsonOutCoin.AssetTag) == 0 {
		assetTag = nil
	} else {
		assetTagInBytes, _, err := base58.Base58Check{}.Decode(jsonOutCoin.AssetTag)
		if err != nil {
			return nil, nil, err
		}
		assetTag, err = new(crypto.Point).FromBytesS(assetTagInBytes)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(jsonOutCoin.Index) == 0 {
		idx = new(big.Int).SetInt64(-1)
	} else {
		idxInBytes, _, err := base58.Base58Check{}.Decode(jsonOutCoin.Index)
		if err != nil {
			return nil, nil, err
		}
		idx = new(big.Int).SetBytes(idxInBytes)
	}

	if jsonOutCoin.Version == "2" {
		coinV2 := new(coin.CoinV2).Init()
		if len(jsonOutCoin.CoinDetailsEncrypted) != 0 {
			coinDetailEncryptedInBytes, _, err := base58.Base58Check{}.Decode(jsonOutCoin.CoinDetailsEncrypted)
			if err != nil {
				return nil, nil, err
			}
			amountEncrypted := new(crypto.Scalar).FromBytesS(coinDetailEncryptedInBytes)
			coinV2.SetAmount(amountEncrypted)
		} else {
			coinV2.SetValue(value)
		}

		coinV2.SetRandomness(randomness)
		coinV2.SetPublicKey(pubkey)
		coinV2.SetCommitment(cm)
		coinV2.SetKeyImage(keyImage)
		coinV2.SetInfo(info)
		coinV2.SetAssetTag(assetTag)
		coinV2.SetSharedRandom(sharedRandom)
		coinV2.SetSharedConcealRandom(sharedConcealRandom)
		coinV2.SetTxRandom(txRandom)

		return coinV2, idx, nil
	} else { //Default is version 1
		pCoinV1 := new(coin.PlainCoinV1).Init()

		pCoinV1.SetRandomness(randomness)
		pCoinV1.SetPublicKey(pubkey)
		pCoinV1.SetCommitment(cm)
		pCoinV1.SetSNDerivator(snd)
		pCoinV1.SetKeyImage(keyImage)
		pCoinV1.SetInfo(info)
		pCoinV1.SetValue(value)

		if len(jsonOutCoin.CoinDetailsEncrypted) != 0 {
			coinDetailEncryptedInBytes, _, err := base58.Base58Check{}.Decode(jsonOutCoin.CoinDetailsEncrypted)
			if err != nil {
				return nil, nil, err
			}

			if len(coinDetailEncryptedInBytes) > 0 {
				coinDetailEncrypted = new(privacy.HybridCipherText)
				err = coinDetailEncrypted.SetBytes(coinDetailEncryptedInBytes)
				if err != nil {
					return nil, nil, err
				}

				coinV1 := new(coin.CoinV1).Init()
				coinV1.CoinDetails = pCoinV1
				coinV1.CoinDetailsEncrypted = coinDetailEncrypted

				return coinV1, idx, nil
			}
		}

		return pCoinV1, idx, nil
	}
}
