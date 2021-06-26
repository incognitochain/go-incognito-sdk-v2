package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

type ConvertedPrice struct {
	FromTokenIDStr string
	ToTokenIDStr   string
	Amount         uint64
	Price          uint64
}

func (server *RPCServer) PDEContributePRV(privKeyStr string, amount string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}
	
	keyWallet, _ := wallet.Base58CheckDeserialize(privKeyStr)
	keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	paymentAddStr := keyWallet.Base58CheckSerialize(wallet.PaymentAddressType)

	//Attempt to contribute to pdex using old address
	paymentAddStr, _ = wallet.GetPaymentAddressV1(paymentAddStr, false)
	query := fmt.Sprintf(`{
				"id": 1,
				"jsonrpc": "1.0",
				"method": "createandsendtxwithprvcontribution",
				"params": [
					"%s",
					{
						"12RxahVABnAVCGP3LGwCn8jkQxgw7z1x14wztHzn455TTVpi1wBq9YGwkRMQg3J4e657AbAnCvYCJSdA9czBUNuCKwGSRQt55Xwz8WA": %s
					},
					-1,
					0,
					{
						"PDEContributionPairID": "newpair",
						"ContributorAddressStr": "%s",
						"ContributedAmount": "%s",
						"TokenIDStr": "0000000000000000000000000000000000000000000000000000000000000004"
					}
				]
			}`, privKeyStr, amount, paymentAddStr, amount)
	return server.SendPostRequestWithQuery(query)
}

func (server *RPCServer) PDEContributeToken(privKeyStr, tokenID, amount string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}
	
	keyWallet, _ := wallet.Base58CheckDeserialize(privKeyStr)
	keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	paymentAddStr := keyWallet.Base58CheckSerialize(wallet.PaymentAddressType)

	//Attempt to contribute to pdex using old address
	paymentAddStr, _ = wallet.GetPaymentAddressV1(paymentAddStr, false)
	query := fmt.Sprintf(`{
				"id": 1,
				"jsonrpc": "1.0",
				"method": "createandsendtxwithptokencontribution",
				"params": [
					"%s",
					{},
					-1,
					0,
					{
						"Privacy": true,
						"TokenID": "%s",
						"TokenTxType": 1,
						"TokenName": "",
						"TokenSymbol": "",
						"TokenAmount": %s,
						"TokenReceivers": {
							"12RxahVABnAVCGP3LGwCn8jkQxgw7z1x14wztHzn455TTVpi1wBq9YGwkRMQg3J4e657AbAnCvYCJSdA9czBUNuCKwGSRQt55Xwz8WA": %s
						},
						"TokenFee": 0,
						"PDEContributionPairID": "newpair",
						"ContributorAddressStr": "%s",
						"ContributedAmount": "%s",
						"TokenIDStr": "%s"
					},
					"",
					0
				]
			}`, privKeyStr, tokenID, amount, amount, paymentAddStr, amount, tokenID)
	return server.SendPostRequestWithQuery(query)
}

func (server *RPCServer) PDEWithdrawContribution(privKeyStr, tokenID1, tokenID2, amountShare string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}
	
	keyWallet, _ := wallet.Base58CheckDeserialize(privKeyStr)
	keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	paymentAddStr := keyWallet.Base58CheckSerialize(wallet.PaymentAddressType) //Attempt to withdraw contribution using new a payment address

	//Attempt to withdraw to pdex using old address
	paymentAddStr, _ = wallet.GetPaymentAddressV1(paymentAddStr, false)
	query := fmt.Sprintf(`{
			"id": 1,
			"jsonrpc": "1.0",
			"method": "createandsendtxwithwithdrawalreq",
		   "params": [
				"%s",
				{
					"12RxahVABnAVCGP3LGwCn8jkQxgw7z1x14wztHzn455TTVpi1wBq9YGwkRMQg3J4e657AbAnCvYCJSdA9czBUNuCKwGSRQt55Xwz8WA": 0
				},
				5,
				-1,
				{
					"WithdrawalShareAmt": %s,
					"WithdrawalToken1IDStr": "%s",
					"WithdrawalToken2IDStr": "%s",
					"WithdrawerAddressStr": "%s"
				}
			]
		}`, privKeyStr, amountShare, tokenID1, tokenID2, paymentAddStr)
	return server.SendPostRequestWithQuery(query)
}

func (server *RPCServer) PDEFeeWithdraw(privKeyStr, tokenID1, tokenID2, amountShare string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}
	
	keyWallet, _ := wallet.Base58CheckDeserialize(privKeyStr)
	keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	paymentAddStr := keyWallet.Base58CheckSerialize(wallet.PaymentAddressType) //Attempt to withdraw contribution using new a payment address
	query := fmt.Sprintf(`{
			"id": 1,
			"jsonrpc": "1.0",
			"method": "createandsendtxwithpdefeewithdrawalreq",
		   "params": [
				"%s",
				{
					"12RxahVABnAVCGP3LGwCn8jkQxgw7z1x14wztHzn455TTVpi1wBq9YGwkRMQg3J4e657AbAnCvYCJSdA9czBUNuCKwGSRQt55Xwz8WA": 0
				},
				5,
				-1,
				{
					"WithdrawalFeeAmt": %s,
					"WithdrawalToken1IDStr": "%s",
					"WithdrawalToken2IDStr": "%s",
					"WithdrawerAddressStr": "%s"
				}
			]
		}`, privKeyStr, amountShare, tokenID1, tokenID2, paymentAddStr)
	return server.SendPostRequestWithQuery(query)
}

func (server *RPCServer) PDETradePRV(privKeyStr, receiverToken, amount string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}
	
	keyWallet, _ := wallet.Base58CheckDeserialize(privKeyStr)
	keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	paymentAddStr := keyWallet.Base58CheckSerialize(wallet.PaymentAddressType)
	query := fmt.Sprintf(`{
			"id": 1,
			"jsonrpc": "1.0",
			"method": "createandsendtxwithprvtradereq",
			"params": [
				"%s",
				{
					"12RxahVABnAVCGP3LGwCn8jkQxgw7z1x14wztHzn455TTVpi1wBq9YGwkRMQg3J4e657AbAnCvYCJSdA9czBUNuCKwGSRQt55Xwz8WA": %s
				},
				-1,
				-1,
				{
					"TokenIDToBuyStr": "%s",
					"TokenIDToSellStr": "0000000000000000000000000000000000000000000000000000000000000004",
					"SellAmount": %s,
					"MinAcceptableAmount": 0,
					"TradingFee": 0,
					"TraderAddressStr": "%s"
				}
			]
		}`, privKeyStr, amount, receiverToken, amount, paymentAddStr)
	return server.SendPostRequestWithQuery(query)
}

func (server *RPCServer) PDETradeToken(privKeyStr, sellToken, amount string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}
	
	keyWallet, _ := wallet.Base58CheckDeserialize(privKeyStr)
	keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	paymentAddStr := keyWallet.Base58CheckSerialize(wallet.PaymentAddressType)
	query := fmt.Sprintf(`{
			"id": 1,
			"jsonrpc": "1.0",
			"method": "createandsendtxwithptokencrosspooltradereq",
			"params": [
				"%s",
				{"12RxahVABnAVCGP3LGwCn8jkQxgw7z1x14wztHzn455TTVpi1wBq9YGwkRMQg3J4e657AbAnCvYCJSdA9czBUNuCKwGSRQt55Xwz8WA": 20},
				-1,
				0,
				{
					"Privacy": true,
					"TokenID": "%s",
					"TokenTxType": 1,
					"TokenName": "",
					"TokenSymbol": "",
					"TokenAmount": 15,
					"TokenReceivers": {
						"12RxahVABnAVCGP3LGwCn8jkQxgw7z1x14wztHzn455TTVpi1wBq9YGwkRMQg3J4e657AbAnCvYCJSdA9czBUNuCKwGSRQt55Xwz8WA": 10
					},
					"TokenFee": 0,
					"TokenIDToBuyStr": "0000000000000000000000000000000000000000000000000000000000000004",
					"TokenIDToSellStr": "%s",
					"SellAmount": 10,
					"MinAcceptableAmount":99999999,
					"TradingFee":20,
					"TraderAddressStr": "%s"
				},
				"",
				0
			]
		}`, privKeyStr, sellToken, sellToken, paymentAddStr)
	return server.SendPostRequestWithQuery(query)
}

func (server *RPCServer) CheckTradeStatus(txHash string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}
	
	method := getPDETradeStatus
	mapParam := make(map[string]interface{})
	mapParam["TxRequestIDStr"] = txHash

	params := make([]interface{}, 0)
	params = append(params, mapParam)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) GetPDEState(beaconHeight uint64) ([]byte, error){
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}
	
	query := fmt.Sprintf(`{
    "id": 1,
    "jsonrpc": "1.0",
    "method": "getpdestate",
    "params": [
        {
            "BeaconHeight": %d
        }
        
    ]
	}`, beaconHeight)

	return server.SendPostRequestWithQuery(query)
}

func (server *RPCServer) ConvertPDEPrice(tokenToSell, tokenToBuy string, amount uint64) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}
	
	method := convertPDEPrices
	mapParam := make(map[string]interface{})
	mapParam["FromTokenIDStr"] = tokenToSell
	mapParam["ToTokenIDStr"] = tokenToBuy
	mapParam["Amount"] = amount

	params := make([]interface{}, 0)
	params = append(params, mapParam)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) GetAllTradesInMempool() ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}
	
	method := getAllTradesInMempool

	params := make([]interface{}, 0)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) GetAllTradesByAddress(addr string) ([]byte, error) {
	method := getAllTradesByAddress
	params := make([]interface{}, 0)

	mapParams := make(map[string]interface{})
	mapParams["PaymentAddress"] = addr
	params = append(params, mapParams)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}
