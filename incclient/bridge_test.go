package incclient

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/eth_bridge/erc20"
	"github.com/incognitochain/go-incognito-sdk-v2/eth_bridge/vault"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"testing"
	"time"
)

//UTILS
const (
	SENDER_KEY         = "0082f4854184d00db3f87d824080460030631a9e4ad7191b93c850e6bc1bdf2e"
	RECEVIER_KEY       = "b446151522b8f1c9d27cacedce93f398a016f84337c1b79fc54c8436af5f7900"
	DEFAULT_ETH_CLIENT = "https://kovan.infura.io/v3/93fe721349134964aa71071a713c5cef"
	ETH_UNIT           = 1e18
	TOKEN_UNIT         = 1e9
	TOKEN_ADDRESS      = "4f96fe3b7a6cf9725f59d353f723c1bdb64ca6aa"
	P_TOKEN_ID         = "c7545459764224a000a9b323850648acf271186238210ce474b505cd17cc93a0"
	P_ETH_ID           = "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854"
	VAULT_ADDRESS      = "0x2f6F03F1b43Eab22f7952bd617A24AB46E970dF7"
)

var ETHClient *ethclient.Client
var ic *IncClient

func InitClients() error {
	var err error

	ETHClient, err = ethclient.Dial(DEFAULT_ETH_CLIENT)
	if err != nil {
		return fmt.Errorf("dial ETHClient error: %v", err)
	}

	ic, err = NewTestNet1Client()
	if err != nil {
		return fmt.Errorf("cannot init new incognito client")
	}

	return nil
}

type Account struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    common.Address
}

func Pad(data string, length int) string {
	if len(data) >= length {
		return data
	}

	for len(data) < length {
		data = "0" + data
	}

	return data
}
func NewETHAccount(hexPrivateKey string) (*Account, error) {
	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot decode hex private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &Account{
		privateKey: privateKey,
		publicKey:  publicKeyECDSA,
		address:    address,
	}, nil

}
func (acc Account) DepositETH(vaultAddress common.Address, incPaymentAddrStr string, gasLimit, gasPrice, depositedAmount uint64) (*common.Hash, error) {
	var err error
	if ETHClient == nil {
		err = InitClients()
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("connected to %v\n", DEFAULT_ETH_CLIENT)

	v, err := vault.NewVault(vaultAddress, ETHClient)
	if err != nil {
		return nil, fmt.Errorf("create new vault error: %v", err)
	}

	//calculate gas price
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = ETHClient.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, fmt.Errorf("cannot get gasPriceBigInt price")
		}
	} else {
		gasPriceBigInt = new(big.Int).SetUint64(gasPrice)
	}

	//calculate gas limit
	if gasLimit == 0 {
		gasLimit, err = ETHClient.EstimateGas(context.Background(), ethereum.CallMsg{To: &vaultAddress, Data: []byte{}})
		if err != nil {
			return nil, fmt.Errorf("estimate gas error: %v", err)
		}
	}

	nonce, err := ETHClient.PendingNonceAt(context.Background(), acc.address)
	if err != nil {
		return nil, fmt.Errorf("get pending nonce error: %v", err)
	}

	auth := bind.NewKeyedTransactor(acc.privateKey)
	auth.Value = big.NewInt(int64(depositedAmount))
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPriceBigInt
	auth.Nonce = new(big.Int).SetUint64(nonce)

	tx, err := v.Deposit(auth, incPaymentAddrStr)
	if err != nil {
		return nil, err
	}
	txHash := tx.Hash()

	if err := Wait(txHash); err != nil {
		return nil, err
	}

	return &txHash, nil
}
func (acc Account) UnShield(vaultAddress common.Address, proof *BurnProof, gasLimit, gasPrice uint64) (*common.Hash, error) {
	var err error
	if ETHClient == nil {
		err = InitClients()
		if err != nil {
			return nil, err
		}
	}

	v, err := vault.NewVault(vaultAddress, ETHClient)
	if err != nil {
		return nil, fmt.Errorf("create new vault error: %v", err)
	}

	auth, err := acc.NewTransactionOpts(vaultAddress, gasPrice, gasLimit, 0, []byte{})
	tx, err := v.Withdraw(auth,
		proof.Instruction,
		proof.Heights[0],
		proof.InstPaths[0],
		proof.InstPathIsLefts[0],
		proof.InstRoots[0],
		proof.BlkData[0],
		proof.SigIdxs[0],
		proof.SigVs[0],
		proof.SigRs[0],
		proof.SigSs[0])

	if err != nil {
		return nil, err
	}

	txHash := tx.Hash()
	return &txHash, nil
}
func (acc Account) DepositERC20(vaultAddress, tokenAddress common.Address, incPaymentAddrStr string, gasLimit, gasPrice, depositedAmount uint64) (*common.Hash, error) {
	var err error
	if ETHClient == nil {
		err = InitClients()
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("Create the approval transaction...\n")

	//Create the ERC20-spending approval transaction
	approvedTx, err := acc.ApproveERC20(vaultAddress, tokenAddress, depositedAmount, 1000000, 0)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Approve transaction: %v\n", approvedTx.Hash().String())
	fmt.Printf("Waiting for the approval tx to be confirmed...\n")
	time.Sleep(20 * time.Second)

	status, err := ic.GetEVMTransactionStatus(approvedTx.Hash().String())
	if err != nil {
		return nil, err
	}

	if status != 1 {
		return nil, fmt.Errorf("approval transaction %v FAILED", approvedTx.Hash().String())
	}

	fmt.Printf("Approval transaction success. Start deposit ERC20 to vault...\n")

	v, err := vault.NewVault(vaultAddress, ETHClient)
	if err != nil {
		return nil, fmt.Errorf("create new vault error: %v", err)
	}

	//calculate gas price
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = ETHClient.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, fmt.Errorf("cannot get gasPriceBigInt price")
		}
	} else {
		gasPriceBigInt = new(big.Int).SetUint64(gasPrice)
	}

	//calculate gas limit
	if gasLimit == 0 {
		gasLimit, err = ETHClient.EstimateGas(context.Background(), ethereum.CallMsg{To: &vaultAddress, Data: []byte{}})
		if err != nil {
			return nil, fmt.Errorf("estimate gas error: %v", err)
		}
	}

	nonce, err := ETHClient.PendingNonceAt(context.Background(), acc.address)
	if err != nil {
		return nil, fmt.Errorf("get pending nonce error: %v", err)
	}

	auth := bind.NewKeyedTransactor(acc.privateKey)
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPriceBigInt
	auth.Nonce = new(big.Int).SetUint64(nonce)

	amount := new(big.Int).SetUint64(depositedAmount)

	tx, err := v.DepositERC20(auth, tokenAddress, amount, incPaymentAddrStr)
	if err != nil {
		return nil, err
	}
	txHash := tx.Hash()

	fmt.Printf("Deposited token successfully!, txHash: %v\n", txHash.String())

	if err := Wait(txHash); err != nil {
		return nil, err
	}

	fmt.Printf("Success!!\n")

	return &txHash, nil
}
func (acc Account) ApproveERC20(vaultAddress, tokenAddress common.Address, approvedAmount, gasLimit, gasPrice uint64) (*types.Transaction, error) {
	var err error
	if ETHClient == nil {
		err = InitClients()
		if err != nil {
			return nil, err
		}
	}

	erc20Token, err := erc20.NewErc20(tokenAddress, ETHClient)
	if err != nil {
		return nil, err
	}

	amount := new(big.Int).SetUint64(approvedAmount)

	//calculate gas price
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = ETHClient.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, fmt.Errorf("cannot get gasPriceBigInt price")
		}
	} else {
		gasPriceBigInt = new(big.Int).SetUint64(gasPrice)
	}

	//calculate gas limit
	if gasLimit == 0 {
		gasLimit, err = ETHClient.EstimateGas(context.Background(), ethereum.CallMsg{To: &vaultAddress, Data: []byte{}})
		if err != nil {
			return nil, fmt.Errorf("estimate gas error: %v", err)
		}
	}

	nonce, err := ETHClient.PendingNonceAt(context.Background(), acc.address)
	if err != nil {
		return nil, fmt.Errorf("get pending nonce error: %v", err)
	}

	auth := bind.NewKeyedTransactor(acc.privateKey)
	auth.GasPrice = gasPriceBigInt
	auth.GasLimit = gasLimit
	auth.Nonce = new(big.Int).SetUint64(nonce)

	txHash, err := erc20Token.Approve(auth, vaultAddress, amount)

	return txHash, err
}
func (acc Account) NewTransactionOpts(destAddr common.Address, gasPrice, gasLimit, amount uint64, data []byte) (*bind.TransactOpts, error) {
	var err error
	if ETHClient == nil {
		err = InitClients()
		if err != nil {
			return nil, err
		}
	}

	//calculate gas limit
	if gasLimit == 0 {
		gasLimit, err = ETHClient.EstimateGas(context.Background(), ethereum.CallMsg{To: &destAddr, Data: data})
		if err != nil {
			return nil, fmt.Errorf("estimate gas error: %v", err)
		}
	}

	nonce, err := ETHClient.PendingNonceAt(context.Background(), acc.address)
	if err != nil {
		return nil, fmt.Errorf("get pending nonce error: %v", err)
	}

	auth := bind.NewKeyedTransactor(acc.privateKey)
	auth.GasPrice = new(big.Int).SetUint64(gasPrice)
	auth.GasLimit = gasLimit
	auth.Nonce = new(big.Int).SetUint64(nonce)
	if amount != 0 {
		value := new(big.Int).SetUint64(amount)
		auth.Value = value
	}

	return auth, nil
}
func GetETHBalance(address common.Address) (uint64, error) {
	var err error
	if ETHClient == nil {
		err = InitClients()
		if err != nil {
			return 0, err
		}
	}

	balance, err := ETHClient.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return 0, fmt.Errorf("get balance error: %v", err)
	}

	pendingBalance, err := ETHClient.PendingBalanceAt(context.Background(), address)
	if err != nil {
		return 0, fmt.Errorf("pending balance error: %v", err)
	}

	balances := make(map[string]uint64)
	balances["balance"] = balance.Uint64()
	balances["pendingBalance"] = pendingBalance.Uint64()

	return balance.Uint64(), nil
}
func GetTokenBalance(address common.Address, scAddress common.Address) (uint64, error) {
	var err error
	if ETHClient == nil {
		err = InitClients()
		if err != nil {
			return 0, err
		}
	}

	instance, err := erc20.NewErc20(scAddress, ETHClient)
	if err != nil {
		return 0, err
	}

	balance, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		return 0, err
	}

	return balance.Uint64(), nil
}

//func GetEthTxStatus(txHash common.Hash) {
//	var err error
//	if ETHClient == nil {
//		err = InitClients()
//		if err != nil {
//			return 0, err
//		}
//	}
//
//	ETHClient.
//}

type BurnProof struct {
	Instruction []byte
	Heights     [2]*big.Int

	InstPaths       [2][][32]byte
	InstPathIsLefts [2][]bool
	InstRoots       [2][32]byte
	BlkData         [2][32]byte
	SigIdxs         [2][]*big.Int
	SigVs           [2][]uint8
	SigRs           [2][][32]byte
	SigSs           [2][][32]byte
}

func DecodeBurnProof(r *jsonresult.GetInstructionProof) (*BurnProof, error) {
	inst := decode(r.Instruction)

	// Block heights
	beaconHeight := big.NewInt(0).SetBytes(decode(r.BeaconHeight))
	bridgeHeight := big.NewInt(0).SetBytes(decode(r.BridgeHeight))
	heights := [2]*big.Int{beaconHeight, bridgeHeight}

	beaconInstRoot := decode32(r.BeaconInstRoot)
	beaconInstPath := make([][32]byte, len(r.BeaconInstPath))
	beaconInstPathIsLeft := make([]bool, len(r.BeaconInstPath))
	for i, path := range r.BeaconInstPath {
		beaconInstPath[i] = decode32(path)
		beaconInstPathIsLeft[i] = r.BeaconInstPathIsLeft[i]
	}
	// fmt.Printf("beaconInstRoot: %x\n", beaconInstRoot)

	beaconBlkData := toByte32(decode(r.BeaconBlkData))

	beaconSigVs, beaconSigRs, beaconSigSs, err := decodeSigs(r.BeaconSigs)
	if err != nil {
		return nil, err
	}

	beaconSigIdxs := []*big.Int{}
	for _, sIdx := range r.BeaconSigIdxs {
		beaconSigIdxs = append(beaconSigIdxs, big.NewInt(int64(sIdx)))
	}

	// For bridge
	bridgeInstRoot := decode32(r.BridgeInstRoot)
	bridgeInstPath := make([][32]byte, len(r.BridgeInstPath))
	bridgeInstPathIsLeft := make([]bool, len(r.BridgeInstPath))
	for i, path := range r.BridgeInstPath {
		bridgeInstPath[i] = decode32(path)
		bridgeInstPathIsLeft[i] = r.BridgeInstPathIsLeft[i]
	}
	// fmt.Printf("bridgeInstRoot: %x\n", bridgeInstRoot)
	bridgeBlkData := toByte32(decode(r.BridgeBlkData))

	bridgeSigVs, bridgeSigRs, bridgeSigSs, err := decodeSigs(r.BridgeSigs)
	if err != nil {
		return nil, err
	}

	bridgeSigIdxs := []*big.Int{}
	for _, sIdx := range r.BridgeSigIdxs {
		bridgeSigIdxs = append(bridgeSigIdxs, big.NewInt(int64(sIdx)))
		// fmt.Printf("bridgeSigIdxs[%d]: %d\n", i, j)
	}

	// Merge beacon and bridge proof
	instPaths := [2][][32]byte{beaconInstPath, bridgeInstPath}
	instPathIsLefts := [2][]bool{beaconInstPathIsLeft, bridgeInstPathIsLeft}
	instRoots := [2][32]byte{beaconInstRoot, bridgeInstRoot}
	blkData := [2][32]byte{beaconBlkData, bridgeBlkData}
	sigIdxs := [2][]*big.Int{beaconSigIdxs, bridgeSigIdxs}
	sigVs := [2][]uint8{beaconSigVs, bridgeSigVs}
	sigRs := [2][][32]byte{beaconSigRs, bridgeSigRs}
	sigSs := [2][][32]byte{beaconSigSs, bridgeSigSs}

	return &BurnProof{
		Instruction:     inst,
		Heights:         heights,
		InstPaths:       instPaths,
		InstPathIsLefts: instPathIsLefts,
		InstRoots:       instRoots,
		BlkData:         blkData,
		SigIdxs:         sigIdxs,
		SigVs:           sigVs,
		SigRs:           sigRs,
		SigSs:           sigSs,
	}, nil
}
func decodeSigs(sigs []string) (sigVs []uint8, sigRs [][32]byte, sigSs [][32]byte, err error) {
	sigVs = make([]uint8, len(sigs))
	sigRs = make([][32]byte, len(sigs))
	sigSs = make([][32]byte, len(sigs))
	for i, sig := range sigs {
		v, r, s, e := key.DecodeECDSASig(sig)
		if e != nil {
			err = e
			return
		}
		sigVs[i] = uint8(v)
		copy(sigRs[i][:], r)
		copy(sigSs[i][:], s)
	}
	return
}
func toByte32(s []byte) [32]byte {
	a := [32]byte{}
	copy(a[:], s)
	return a
}
func decode(s string) []byte {
	d, _ := hex.DecodeString(s)
	return d
}
func decode32(s string) [32]byte {
	return toByte32(decode(s))
}
func keccak256(b ...[]byte) [32]byte {
	h := crypto.Keccak256(b...)
	r := [32]byte{}
	copy(r[:], h)
	return r
}

func Wait(tx common.Hash) error {
	if ETHClient == nil {
		err := InitClients()
		if err != nil {
			return err
		}
	}

	for range time.Tick(10 * time.Second) {
		_, err := ETHClient.TransactionReceipt(context.Background(), tx)
		if err == nil {
			break
		} else if err == ethereum.NotFound {
			continue
		} else {
			return err
		}
	}
	return nil
}

//END UTILS

//TEST FUNCTIONS
func TestIncClient_ShieldETH(t *testing.T) {
	err := InitClients()
	if err != nil {
		panic(err)
	}

	//Init an Ethereum account
	acc, err := NewETHAccount(SENDER_KEY)
	if err != nil {
		panic(err)
	}

	//Incognito keys
	privateKey := "112t8rnZDRztVgPjbYQiXS7mJgaTzn66NvHD7Vus2SrhSAY611AzADsPFzKjKQCKWTgbkgYrCPo9atvSMoCf9KT23Sc7Js9RKhzbNJkxpJU6"
	incAddr := PrivateKeyToPaymentAddress(privateKey, 1)

	//incognito vault contract address
	vaultAddress := common.HexToAddress(VAULT_ADDRESS)

	oldBalance, err := ic.GetBalance(privateKey, P_ETH_ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current balance of token %v: %v\n", P_ETH_ID, oldBalance)

	gasPrice := uint64(50 * 1e9)
	gasLimit := uint64(100000)
	ethAmount := uint64(0.05 * ETH_UNIT)
	ethTxHash, err := acc.DepositETH(vaultAddress, incAddr, gasLimit, gasPrice, ethAmount)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Deposit transaction: %v\n", ethTxHash.String())
	fmt.Printf("Waiting for 15 confirmations...\n")
	time.Sleep(200 * time.Second)

	fmt.Printf("Start shielding eth...")

	ethProof, pETHAmount, err := ic.GetEVMDepositProof(ethTxHash.String())
	if err != nil {
		panic(err)
	}

	txHashStr, err := ic.CreateAndSendIssuingETHRequestTransaction(privateKey, P_ETH_ID, *ethProof)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Finish shielding: %v\n", txHashStr)
	time.Sleep(10 * time.Second)

	fmt.Printf("Check shielding status\n")
	for {
		status, err := ic.CheckShieldStatus(txHashStr)
		if err != nil {
			panic(err)
		}
		if status == 1 || status == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		if status == 2 {
			fmt.Printf("Shielding accepted, start checking balance\n")
			break
		} else {
			panic("Shield rejected!")
		}
	}

	for {
		newBalance, err := ic.GetBalance(privateKey, P_ETH_ID)
		if err != nil {
			panic(err)
		}
		updatedAmount := newBalance - oldBalance
		if updatedAmount != 0 {
			if updatedAmount != pETHAmount {
				panic(fmt.Sprintf("expected %v, got %v\n", pETHAmount, updatedAmount))
			}
			fmt.Printf("Balance updated!\nnewBalance %v, increasedAmount %v, ethAmount %v\n", newBalance, updatedAmount, pETHAmount)
			break
		}
		fmt.Printf("Balance not updated, sleeping for more...\n")
		time.Sleep(5 * time.Second)
	}
}

func TestIncClient_ShieldERC20(t *testing.T) {
	err := InitClients()
	if err != nil {
		panic(err)
	}

	//Init an Ethereum account
	acc, err := NewETHAccount(SENDER_KEY)
	if err != nil {
		panic(err)
	}

	tokenAddress := common.HexToAddress(TOKEN_ADDRESS)
	oldTokenBalance, err := GetTokenBalance(acc.address, tokenAddress)
	if err != nil {
		panic(err)
	}

	oldETHBalance, err := GetETHBalance(acc.address)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Remote account: %v, tokenAddress: %v, balance %v, eth %v\n", acc.address.String(), TOKEN_ADDRESS, oldTokenBalance, oldETHBalance)

	//Incognito keys
	privateKey := "112t8rnZDRztVgPjbYQiXS7mJgaTzn66NvHD7Vus2SrhSAY611AzADsPFzKjKQCKWTgbkgYrCPo9atvSMoCf9KT23Sc7Js9RKhzbNJkxpJU6"
	incAddr := PrivateKeyToPaymentAddress(privateKey, 1)

	//incognito vault contract address
	vaultAddress := common.HexToAddress(VAULT_ADDRESS)

	oldBalance, err := ic.GetBalance(privateKey, P_TOKEN_ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current balance of token %v: %v\n", P_TOKEN_ID, oldBalance)

	gasPrice := uint64(50 * 1e9)
	gasLimit := uint64(1000000)
	depositAmount := uint64(0.1 * ETH_UNIT)
	fmt.Printf("Deposit amount: %v\n", depositAmount)

	ethTxHash, err := acc.DepositERC20(vaultAddress, tokenAddress, incAddr, gasLimit, gasPrice, depositAmount)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Deposited %v, token %v, txHash: %v\n", depositAmount, TOKEN_ADDRESS, ethTxHash.String())
	fmt.Printf("Waiting for 15 confirmations...\n")
	time.Sleep(90 * time.Second)

	status, err := ic.GetEVMTransactionStatus(ethTxHash.String())
	if err != nil {
		panic(err)
	}
	if status != 1 {
		panic(fmt.Sprintf("transaction %v failed on the ETH network", ethTxHash.String()))
	}

	fmt.Printf("Start shielding token...\n")

	//ethTxHash := common.HexToHash("0x75807ac2052d7a612d857aa4e5e3fea3e0a07007543c9fb7b5bea1d6cba069f4")

	ethProof, _, err := ic.GetEVMDepositProof(ethTxHash.String())
	if err != nil {
		panic(err)
	}

	txHashStr, err := ic.CreateAndSendIssuingETHRequestTransaction(privateKey, P_TOKEN_ID, *ethProof)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Finish shielding: %v\n", txHashStr)
	time.Sleep(10 * time.Second)

	fmt.Printf("Check shielding status\n")
	for {
		status, err := ic.CheckShieldStatus(txHashStr)
		if err != nil {
			panic(err)
		}
		if status == 1 || status == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		if status == 2 {
			fmt.Printf("Shielding accepted, start checking balance\n")
			break
		} else {
			panic(fmt.Sprintf("Shield rejected, status: %v\n", status))
		}
	}

	for {
		newBalance, err := ic.GetBalance(privateKey, P_TOKEN_ID)
		if err != nil {
			panic(err)
		}
		updatedAmount := newBalance - oldBalance
		if updatedAmount != 0 {
			expectedIncrease := depositAmount / 1e9
			if updatedAmount != expectedIncrease {
				panic(fmt.Sprintf("expected %v, got %v\n", expectedIncrease, updatedAmount))
			}
			fmt.Printf("Balance updated!\nnewBalance %v, increasedAmount %v, depositAmount %v\n", newBalance, updatedAmount, expectedIncrease)

			newETHBalance, err := GetETHBalance(acc.address)
			if err == nil {
				changedAmount := oldETHBalance - newETHBalance
				fmt.Printf("NewETHbalance: %v, changedAmount: %v\n", newETHBalance, changedAmount)
			}
			break
		}
		fmt.Printf("Balance not updated, sleeping for more...\n")
		time.Sleep(5 * time.Second)
	}
}

func TestIncClient_UnShieldETH(t *testing.T) {
	err := InitClients()
	if err != nil {
		panic(err)
	}

	//Init an Ethereum account
	acc, err := NewETHAccount(SENDER_KEY)
	if err != nil {
		panic(err)
	}

	oldBalance, err := GetETHBalance(acc.address)
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rnZDRztVgPjbYQiXS7mJgaTzn66NvHD7Vus2SrhSAY611AzADsPFzKjKQCKWTgbkgYrCPo9atvSMoCf9KT23Sc7Js9RKhzbNJkxpJU6"
	remoteAddr := acc.address.String()
	burnedAmount := uint64(50000000)

	fmt.Printf("Remote account: %v, balance %v\n", remoteAddr, oldBalance)

	burnedTxHash, err := ic.CreateAndSendBurningRequestTransaction(privateKey, remoteAddr, P_ETH_ID, burnedAmount)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wait for tx %v to be confirmed\n", burnedTxHash)
	for {
		isInBlock, err := ic.CheckTxInBlock(burnedTxHash)
		if err != nil {
			panic(err)
		}

		if !isInBlock {
			fmt.Printf("Tx %v is currently in mempool\n", burnedTxHash)
			time.Sleep(10 * time.Second)
		} else {
			fmt.Printf("Tx %v is in block\n", burnedTxHash)
			fmt.Printf("Sleep 40 seconds for getting burning proof\n")
			time.Sleep(40 * time.Second)
			break
		}
	}

	fmt.Printf("Start to retrieve the burning proof\n")
	//burnedTxHash := "a4c6585de955e29707ecf1658e268c6182513b3d0fd765d04b630e3866a18c19"
	burningProofResult, err := ic.GetBurnProof(burnedTxHash)
	if err != nil {
		panic(err)
	}

	burnProof, err := DecodeBurnProof(burningProofResult)
	if err != nil {
		panic(err)
	}

	jsb, _ := json.Marshal(burningProofResult)

	fmt.Printf("Burn proof from Incog: %v\n", string(jsb))

	fmt.Printf("Finish getting the burning proof\n")

	fmt.Printf("Start submitting the proof\n")

	vaultAddress := common.HexToAddress(VAULT_ADDRESS)
	gasPrice := uint64(50 * 1e9)
	gasLimit := uint64(1000000)
	txHash, err := acc.UnShield(vaultAddress, burnProof, gasLimit, gasPrice)
	if err != nil {
		panic(err)
	}

	if err := Wait(*txHash); err != nil {
		panic(err)
	}
	fmt.Printf("Unshield tx: %v\n", txHash.String())

	fmt.Printf("Check unshield tx status...\n")
	receipt, err := ic.GetEVMTxReceipt(txHash.String())
	if err != nil {
		panic(err)
	}

	if receipt.Status != 1 {
		panic(fmt.Errorf("unshield tx FAILED\n"))
	}

	fmt.Printf("Check balance updated for %v\n", acc.address.String())
	for {
		newBalance, err := GetETHBalance(acc.address)
		if err != nil {
			panic(err)
		}
		if newBalance < oldBalance {
			time.Sleep(10 * time.Second)
			continue
		}
		updatedAmount := newBalance - oldBalance
		if updatedAmount != 0 {
			fmt.Printf("Balance updated!\nnewBalance: %v, increasedAmount %v\n", newBalance, updatedAmount)
			break
		} else {
			fmt.Printf("Balance not updated, sleeping for more...\n")
			time.Sleep(5 * time.Second)
		}
	}
}

func TestIncClient_UnShieldERC20(t *testing.T) {
	err := InitClients()
	if err != nil {
		panic(err)
	}

	//Init an Ethereum account
	acc, err := NewETHAccount(SENDER_KEY)
	if err != nil {
		panic(err)
	}

	//Check balance of the ERC20 token before unshield
	tokenAddress := common.HexToAddress(TOKEN_ADDRESS)
	oldBalance, err := GetTokenBalance(acc.address, tokenAddress)
	if err != nil {
		panic(err)
	}

	oldETHBalance, err := GetETHBalance(acc.address)
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rnZDRztVgPjbYQiXS7mJgaTzn66NvHD7Vus2SrhSAY611AzADsPFzKjKQCKWTgbkgYrCPo9atvSMoCf9KT23Sc7Js9RKhzbNJkxpJU6"
	tokenIDStr := "c7545459764224a000a9b323850648acf271186238210ce474b505cd17cc93a0" //incognito tokenID for pDAI
	remoteAddr := acc.address.String()
	burnedAmount := uint64(10000000)

	fmt.Printf("Remote account: %v, tokenAddress: %v, balance %v, eth %v\n", remoteAddr, TOKEN_ADDRESS, oldBalance, oldETHBalance)
	fmt.Printf("Unshield amount: %v\n", burnedAmount)
	burnedTxHash, err := ic.CreateAndSendBurningRequestTransaction(privateKey, remoteAddr, tokenIDStr, burnedAmount)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wait for tx %v to be confirmed\n", burnedTxHash)
	for {
		isInBlock, err := ic.CheckTxInBlock(burnedTxHash)
		if err != nil {
			panic(err)
		}

		if !isInBlock {
			fmt.Printf("Tx %v is currently in mempool\n", burnedTxHash)
			time.Sleep(10 * time.Second)
		} else {
			fmt.Printf("Tx %v is in block\n", burnedTxHash)
			fmt.Printf("Sleep 40 seconds for getting burning proof\n")
			time.Sleep(40 * time.Second)
			break
		}
	}

	fmt.Printf("Start to retrieve the burning proof\n")
	//burnedTxHash := "a4c6585de955e29707ecf1658e268c6182513b3d0fd765d04b630e3866a18c19"
	burningProofResult, err := ic.GetBurnProof(burnedTxHash)
	if err != nil {
		panic(err)
	}

	burnProof, err := DecodeBurnProof(burningProofResult)
	if err != nil {
		panic(err)
	}

	jsb, _ := json.Marshal(burningProofResult)

	fmt.Printf("Burn proof from Incog: %v\n", string(jsb))

	fmt.Printf("Finish getting the burning proof\n")

	fmt.Printf("Start submitting the proof\n")

	vaultAddress := common.HexToAddress(VAULT_ADDRESS)
	gasPrice := uint64(50 * 1e9)
	gasLimit := uint64(1000000)
	txHash, err := acc.UnShield(vaultAddress, burnProof, gasLimit, gasPrice)
	if err != nil {
		panic(err)
	}

	if err := Wait(*txHash); err != nil {
		panic(err)
	}
	fmt.Printf("Unshield tx: %v\n", txHash.String())

	fmt.Printf("Check unshield tx status...\n")
	receipt, err := ic.GetEVMTxReceipt(txHash.String())
	if err != nil {
		panic(err)
	}

	if receipt.Status != 1 {
		panic(fmt.Errorf("unshield tx FAILED\n"))
	}

	fmt.Printf("Check balance updated for %v\n", acc.address.String())
	for {
		newBalance, err := GetTokenBalance(acc.address, tokenAddress)
		if err != nil {
			panic(err)
		}
		if newBalance < oldBalance {
			time.Sleep(10 * time.Second)
			continue
		}
		updatedAmount := newBalance - oldBalance
		if updatedAmount != 0 {
			if updatedAmount != burnedAmount*1e9 {
				panic(fmt.Sprintf("expected received amount: %v, got %v", burnedAmount*1e9, updatedAmount))
			}
			fmt.Printf("Balance updated!\nnewBalance: %v, increasedAmount %v\n", newBalance, updatedAmount)

			newETHBalance, err := GetETHBalance(acc.address)
			if err == nil {
				changedAmount := oldETHBalance - newETHBalance
				fmt.Printf("NewETHbalance: %v, changedAmount: %v\n", newETHBalance, changedAmount)
			}
			break
		} else {
			fmt.Printf("Balance not updated, sleeping for more...\n")
			time.Sleep(5 * time.Second)
		}
	}
}

func TestIncClient_GetETHTxReceipt(t *testing.T) {
	err := InitClients()
	if err != nil {
		panic(err)
	}

	txHash := "0xc400656111f353ef021f3f65711461679e4e1227071411c2789cac762e8948bb"

	receipt, err := ic.GetEVMTxReceipt(txHash)
	if err != nil {
		panic(err)
	}

	jsb, err := json.Marshal(receipt)
	if err != nil {
		panic(err)
	}

	fmt.Printf(string(jsb))
}

func TestIncClient_GetMostRecentETHBlockNumber(t *testing.T) {
	err := InitClients()
	if err != nil {
		panic(err)
	}

	blockNum, err := ic.GetMostRecentEVMBlockNumber()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current blockNum: %v\n", blockNum)
}

//END TEST FUNCTIONS
