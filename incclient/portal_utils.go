package incclient

import (
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
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

// GeneratePortalShieldingAddress returns the multi-sig shielding address for a given payment address and a tokenID.
func (client *IncClient) GeneratePortalShieldingAddress(paymentAddressStr, tokenIDStr string) (string, error) {
	var res string
	var err error

	if client.btcPortalParams != nil {
		if tokenIDStr != client.btcPortalParams.TokenID {
			return "", fmt.Errorf("tokenID %v not supported by the v4 Portal", tokenIDStr)
		}

		pubKeys := make([][]byte, 0)
		if paymentAddressStr == "" {
			pubKeys = client.btcPortalParams.MasterPubKeys[:]
		} else {
			_, err = AssertPaymentAddressAndTxVersion(paymentAddressStr, 2)
			if err != nil {
				return "", fmt.Errorf("invalid payment address: %v", err)
			}

			chainCode := chainhash.HashB([]byte(paymentAddressStr))
			for idx, masterPubKey := range client.btcPortalParams.MasterPubKeys {
				// generate BTC child public key for this Incognito address
				extendedBTCPublicKey := hdkeychain.NewExtendedKey(client.btcPortalParams.ChainParams.HDPublicKeyID[:], masterPubKey, chainCode, []byte{}, 0, 0, false)
				// extendedBTCChildPubKey, err := extendedBTCPublicKey.ECPubKey()
				// if err != nil {
				// 	return "", err
				// }

				childPubKey, err := extendedBTCPublicKey.ECPubKey()
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

		rpcRes, err := client.generatePortalShieldingAddressFromRPC(paymentAddressStr, tokenIDStr)
		if err != nil {
			return "", err
		}
		Logger.Println("Generated shielding addresses match!!")

		if rpcRes != res {
			return "", fmt.Errorf("rpc result (%v) and client result (%v) mismatch, please double check the v4 Portal configuration", rpcRes, res)
		}
	}

	res, err = client.generatePortalShieldingAddressFromRPC(paymentAddressStr, tokenIDStr)
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

// generatePortalShieldingAddressFromRPC returns the multi-sig shielding address for a given payment address and a tokenID
// via an RPC when using the Portal.
func (client *IncClient) generatePortalShieldingAddressFromRPC(paymentAddressStr, tokenIDStr string) (string, error) {
	responseInBytes, err := client.rpcServer.GenerateShieldingMultiSigAddress(paymentAddressStr, tokenIDStr)
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
