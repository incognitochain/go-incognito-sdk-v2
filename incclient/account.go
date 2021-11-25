package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// GetBalance retrieves the current tokenID balance of a private key.
func (client *IncClient) GetBalance(privateKey, tokenID string) (uint64, error) {
	unspentCoins, _, err := client.GetUnspentOutputCoins(privateKey, tokenID, 0)
	if err != nil {
		return 0, err
	}

	balance := uint64(0)
	for _, unspentCoin := range unspentCoins {
		balance += unspentCoin.GetValue()
	}

	return balance, nil
}

// GetAllBalancesV2 returns all non-zero balances of a private key.
// This function assumes that all v1 output coins have been converted to v1, and only returns the balances calculated with
// v2 coins (except for PRV). In case you still have v1 UTXOs, try using the regular `GetBalance` function.
func (client *IncClient) GetAllBalancesV2(privateKey string) (map[string]uint64, error) {
	res := make(map[string]uint64)
	allUTXOs, _, err := client.GetAllUTXOsV2(privateKey)
	if err != nil {
		return nil, err
	}

	for tokenID, utxoList := range allUTXOs {
		balance := uint64(0)
		for _, utxo := range utxoList {
			balance += utxo.GetValue()
		}
		if balance > 0 {
			res[tokenID] = balance
		}
	}

	return res, nil
}

// GetAllNFTs returns all NFTs belonging to a private key.
func (client *IncClient) GetAllNFTs(privateKey string) ([]string, error) {
	utxoList, _, err := client.GetUnspentOutputCoins(privateKey, common.ConfidentialAssetID.String(), 0)
	if err != nil {
		return nil, err
	}
	if len(utxoList) == 0 {
		return nil, fmt.Errorf("no UTXO found")
	}
	Logger.Printf("#UTXOs: %v\n", len(utxoList))

	allNFTs, err := client.GetListNftIDs(0)
	if err != nil {
		return nil, err
	}
	nftList := make([]string, 0)
	for tokenID, _ := range allNFTs {
		nftList = append(nftList, tokenID)
	}
	Logger.Printf("#Nfts: %v\n", len(allNFTs))

	rawAssetTags, err = BuildAssetTags(nftList)
	if err != nil {
		return nil, err
	}

	w, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	for _, utxo := range utxoList {
		if utxo.GetValue() != 1 {
			continue
		}
		v2Coin, ok := utxo.(*coin.CoinV2)
		if !ok {
			return nil, fmt.Errorf("cannot cast UTXO %v to a CoinV2", base58.Base58Check{}.Encode(utxo.GetPublicKey().ToBytesS(), 0))
		}
		tokenId, _ := v2Coin.GetTokenId(&(w.KeySet), rawAssetTags)
		if tokenId == nil {
			continue
		}
		if _, ok := allNFTs[tokenId.String()]; ok {
			res = append(res, tokenId.String())
		}
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("no NFT found")
	}
	return res, nil
}

// ImportAccount imports a BIP39 mnemonic string and finds all child keys derived from the mnemonic. The first return KeyWallet
// is the master wallet, which is used to derive the rest of child KeyWallet.
// For child KeyWallets, we start with childIdx = 1 and stops at the index when there is no transaction found for the child key.
func (client *IncClient) ImportAccount(mnemonic string) ([]*wallet.KeyWallet, error) {
	res := make([]*wallet.KeyWallet, 0)

	masterWallet, err := wallet.NewMasterKeyFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	res = append(res, masterWallet)

	childIdx := uint32(1)
	for {
		childWallet, err := masterWallet.DeriveChild(childIdx)
		if err != nil {
			return nil, fmt.Errorf("childIdx %v error: %v", childIdx, err)
		}

		addr, _ := childWallet.GetPaymentAddress()
		receivedTxs, err := client.GetTransactionHashesByReceiver(addr)
		if err != nil {
			return nil, fmt.Errorf("childIdx %v error: %v", childIdx, err)
		}

		if len(receivedTxs) > 0 {
			res = append(res, childWallet)
			childIdx++
		} else {
			return res, nil
		}
	}
}
