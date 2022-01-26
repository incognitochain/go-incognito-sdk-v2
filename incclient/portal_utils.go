package incclient

import (
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"math"
	"math/big"
)

// GetPortalShieldingRequestStatus retrieves the status of a port shielding request.
func (client *IncClient) GetPortalShieldingRequestStatus(shieldID string) (*metadata.PortalShieldingRequestStatus, error) {
	responseInBytes, err := client.rpcServer.GetPortalShieldingRequestStatus(shieldID)
	if err != nil {
		return nil, err
	}

	var res *metadata.PortalShieldingRequestStatus
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GeneratePortalShieldingAddress returns the multi-sig shielding address for a given seed and a tokenID.
func (client *IncClient) GeneratePortalShieldingAddress(chainCodeStr, tokenIDStr string) (string, error) {
	var res string
	var err error

	if client.btcPortalParams != nil {
		if tokenIDStr != client.btcPortalParams.TokenID {
			return "", fmt.Errorf("tokenID %v not supported by the v4 Portal", tokenIDStr)
		}

		pubKeys := make([][]byte, 0)
		if chainCodeStr == "" {
			pubKeys = client.btcPortalParams.MasterPubKeys[:]
		} else {
			var chainCode []byte
			_, err = AssertPaymentAddressAndTxVersion(chainCodeStr, 2)
			if err != nil {
				depositPubKey, _, err := base58.Base58Check{}.Decode(chainCodeStr)
				if err != nil || len(depositPubKey) != 32 {
					return "", fmt.Errorf("invalid chain-code")
				}
				chainCode = chainhash.HashB(depositPubKey)
			} else {
				chainhash.HashB([]byte(chainCodeStr))
			}

			for idx, masterPubKey := range client.btcPortalParams.MasterPubKeys {
				// generate BTC child public key for this Incognito address
				extendedBTCPublicKey := hdkeychain.NewExtendedKey(client.btcPortalParams.ChainParams.HDPublicKeyID[:], masterPubKey, chainCode, []byte{}, 0, 0, false)
				extendedBTCChildPubKey, err := extendedBTCPublicKey.Child(0)
				if err != nil {
					return "", err
				}

				childPubKey, err := extendedBTCChildPubKey.ECPubKey()
				if err != nil {
					return "", fmt.Errorf("master BTC Public Key (#%v) %v is invalid - Error %v", idx, masterPubKey, err)
				}
				pubKeys = append(pubKeys, childPubKey.SerializeCompressed())
			}

			// create redeem script for m of n multi-sig
			builder := txscript.NewScriptBuilder()
			// add the minimum number of needed signatures
			builder.AddOp(byte(txscript.OP_1 - 1 + client.btcPortalParams.NumRequiredSigs))
			// add the public key to redeem script
			for _, pubKey := range pubKeys {
				builder.AddData(pubKey)
			}
			// add the total number of public keys in the multi-sig script
			builder.AddOp(byte(txscript.OP_1 - 1 + len(pubKeys)))
			// add the check-multi-sig op-code
			builder.AddOp(txscript.OP_CHECKMULTISIG)

			redeemScript, err := builder.Script()
			if err != nil {
				return "", fmt.Errorf("could not build script - Error %v", err)
			}

			// generate P2WSH address
			scriptHash := sha256.Sum256(redeemScript)
			addr, err := btcutil.NewAddressWitnessScriptHash(scriptHash[:], client.btcPortalParams.ChainParams)
			if err != nil {
				return "", fmt.Errorf("could not generate address from script - Error %v", err)
			}
			res = addr.EncodeAddress()
		}

		// call RPCs to double-check
		rpcRes, err := client.generatePortalShieldingAddressFromRPC(chainCodeStr, tokenIDStr)
		if err != nil {
			return "", err
		}
		Logger.Println("Generated shielding addresses match!!")

		if rpcRes != res {
			return "", fmt.Errorf("rpc result (%v) and client result (%v) mismatch, please double check the v4 Portal configuration", rpcRes, res)
		}
	}

	res, err = client.generatePortalShieldingAddressFromRPC(chainCodeStr, tokenIDStr)
	if err != nil {
		return "", err
	}

	return res, nil
}

// GetPortalUnShieldingRequestStatus retrieves the status of a port un-shielding request.
func (client *IncClient) GetPortalUnShieldingRequestStatus(unShieldID string) (*metadata.PortalUnshieldRequestStatus, error) {
	responseInBytes, err := client.rpcServer.GetPortalUnShieldingRequestStatus(unShieldID)
	if err != nil {
		return nil, err
	}

	var res *metadata.PortalUnshieldRequestStatus
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GenerateDepositKeyFromPrivateKey generates a new OTDepositKey from the given privateKey, tokenID and index.
func (client *IncClient) GenerateDepositKeyFromPrivateKey(privateKeyStr, tokenIDStr string, index uint64) (*key.OTDepositKey, error) {
	w, err := wallet.Base58CheckDeserialize(privateKeyStr)
	if err != nil || len(w.KeySet.PrivateKey[:]) == 0 {
		return nil, fmt.Errorf("invalid privateKey")
	}

	tokenID, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, err
	}

	tmp := append([]byte(common.PortalV4DepositKeyGenSeed), tokenID[:]...)
	masterDepositSeed := common.SHA256(append(w.KeySet.PrivateKey[:], tmp...))
	indexBig := new(big.Int).SetUint64(index)

	privateKey := crypto.HashToScalar(append(masterDepositSeed, indexBig.Bytes()...))
	pubKey := new(crypto.Point).ScalarMultBase(privateKey)

	return &key.OTDepositKey{
		PrivateKey: privateKey.ToBytesS(),
		PublicKey:  pubKey.ToBytesS(),
		Index:      index,
	}, nil
}

// GetNextOTDepositKey returns the next un-used deposit key and its corresponding depositing address.
func (client *IncClient) GetNextOTDepositKey(privateKeyStr, tokenIDStr string) (*key.OTDepositKey, string, error) {
	tmpKey, err := client.GenerateDepositKeyFromPrivateKey(privateKeyStr, tokenIDStr, 0)
	if err != nil {
		return nil, "", err
	}
	tmpPubKeyStr := base58.Base58Check{}.Encode(tmpKey.PublicKey, 0)

	exists, err := client.HasDepositPubKeys([]string{tmpPubKeyStr})
	if err != nil {
		return nil, "", err
	}

	if exists[tmpPubKeyStr] {
		// Perform binary-search for the un-used index
		lower := uint64(1)
		upper := uint64(math.MaxUint64)
		for {
			currentIndex := uint64(math.Max(float64(lower), 1))
			tmpKey, err = client.GenerateDepositKeyFromPrivateKey(privateKeyStr, tokenIDStr, currentIndex)
			if err != nil {
				return nil, "", fmt.Errorf("generating depositKey at index %v error: %v", lower, err)
			}
			tmpPubKeyStr = base58.Base58Check{}.Encode(tmpKey.PublicKey, 0)
			exists, err = client.HasDepositPubKeys([]string{tmpPubKeyStr})
			if err != nil {
				return nil, "", err
			}
			if exists[tmpPubKeyStr] {
				lower = currentIndex
				currentIndex = 2 * lower
			} else {
				upper = currentIndex
				break
			}
		}

		currentIndex := lower
		for lower < upper-1 {
			tmpKey, err = client.GenerateDepositKeyFromPrivateKey(privateKeyStr, tokenIDStr, currentIndex)
			if err != nil {
				return nil, "", fmt.Errorf("generating depositKey at index %v error: %v", lower, err)
			}
			tmpPubKeyStr = base58.Base58Check{}.Encode(tmpKey.PublicKey, 0)
			exists, err = client.HasDepositPubKeys([]string{tmpPubKeyStr})
			if err != nil {
				return nil, "", err
			}

			if !exists[tmpPubKeyStr] {
				upper = currentIndex
			} else {
				lower = currentIndex
			}
			currentIndex = (lower + upper) / 2
		}
	}

	depositAddress, err := client.GeneratePortalShieldingAddress(tmpPubKeyStr, tokenIDStr)
	if err != nil {
		return nil, "", err
	}

	return tmpKey, depositAddress, nil
}

// HasDepositPubKeys checks if one-time deposit keys have been used.
func (client *IncClient) HasDepositPubKeys(depositPubKeys []string) (map[string]bool, error) {
	responseInBytes, err := client.rpcServer.HasOTDepositKey(depositPubKeys)
	if err != nil {
		return nil, err
	}

	res := make(map[string]bool)
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetDepositTxsByPubKeys retrieves the Incognito depositing transactions for a given list of depositing public keys.
func (client *IncClient) GetDepositTxsByPubKeys(depositPubKeys []string) (map[string]string, error) {
	responseInBytes, err := client.rpcServer.GetDepositTxsByPubKeys(depositPubKeys)
	if err != nil {
		return nil, err
	}

	res := make(map[string]string)
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetDepositByOTKeyHistory retrieves the depositing (using OTDepositKey) history of a privateKeyStr for a tokenID.
func (client *IncClient) GetDepositByOTKeyHistory(privateKeyStr, tokenID string) (map[string]*metadata.PortalShieldingRequestStatus, error) {
	depositPubKeys := make([]string, 0)

	nextAvailableDepositKey, _, err := client.GetNextOTDepositKey(privateKeyStr, tokenID)
	if err != nil {
		return nil, err
	}
	for index := uint64(0); index < nextAvailableDepositKey.Index; index++ {
		depositKey, err := client.GenerateDepositKeyFromPrivateKey(privateKeyStr, tokenID, index)
		if err != nil {
			return nil, err
		}
		depositPubKey := base58.Base58Check{}.Encode(depositKey.PublicKey, 0)
		depositPubKeys = append(depositPubKeys, depositPubKey)
	}

	if len(depositPubKeys) == 0 {
		return nil, fmt.Errorf("no deposit history found")
	}

	depositTxs, err := client.GetDepositTxsByPubKeys(depositPubKeys)
	if err != nil {
		return nil, err
	}
	res := make(map[string]*metadata.PortalShieldingRequestStatus)
	for pubKeyStr, txHash := range depositTxs {
		status, err := client.GetPortalShieldingRequestStatus(txHash)
		if err != nil {
			return nil, err
		}
		res[pubKeyStr] = status
	}

	return res, nil
}

// generatePortalShieldingAddressFromRPC returns the multi-sig shielding address for a given payment address and a tokenID
// via an RPC when using the Portal.
func (client *IncClient) generatePortalShieldingAddressFromRPC(paymentAddressStr, tokenIDStr string) (string, error) {
	responseInBytes, err := client.rpcServer.GenerateDepositAddress(paymentAddressStr, tokenIDStr)
	if err != nil {
		return "", err
	}

	var res string
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return "", err
	}

	return res, nil
}

// SignDepositData signs the given depositing data using the given OTDepositKey.
// Data is the raw depositing data (not hashed).
func SignDepositData(depositKey *key.OTDepositKey, data []byte) ([]byte, error) {
	schnorrPrivateKey := new(privacy.SchnorrPrivateKey)
	schnorrPrivateKey.Set(new(crypto.Scalar).FromBytesS(depositKey.PrivateKey), crypto.RandomScalar())

	digestedData := common.HashB(data)
	sig, err := schnorrPrivateKey.Sign(digestedData)
	if err != nil {
		return nil, err
	}

	return sig.Bytes(), nil
}
