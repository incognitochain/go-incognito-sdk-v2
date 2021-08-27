package incclient

import (
	"fmt"
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

//// GetBalanceAll returns all non-zero balances (for all tokenIDs) of a private key.
//func (client *IncClient) GetBalanceAll(privateKey string) (map[string]uint64, error) {
//	res := make(map[string]uint64)
//	prvBalance, err := client.GetBalance(privateKey, common.PRVIDStr)
//	if err != nil {
//		return nil, err
//	}
//	if prvBalance > 0 {
//		res[common.PRVIDStr] = prvBalance
//	}
//
//	tokenList, err := client.GetListToken()
//	if err != nil {
//		return nil, err
//	}
//
//	if client.cache != nil && client.cache.isRunning {
//		tokenCount := 0
//		for tokenID := range tokenList {
//			reSync := tokenCount == 0
//			unspentCoins, _, err := client.GetUnspentOutputCoinsFromCache(privateKey, tokenID, 0, reSync)
//			if err != nil {
//				return nil, err
//			}
//			balance := uint64(0)
//			for _, unspentCoin := range unspentCoins {
//				balance += unspentCoin.GetValue()
//			}
//			if balance > 0 {
//				res[tokenID] = balance
//			}
//			tokenCount++
//		}
//		return res, nil
//	}
//
//	for tokenID := range tokenList {
//		tmpBalance, err := client.GetBalance(privateKey, tokenID)
//		if err != nil {
//			return nil, err
//		}
//		if tmpBalance > 0 {
//			res[tokenID] = tmpBalance
//		}
//	}
//
//	return res, nil
//}

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
