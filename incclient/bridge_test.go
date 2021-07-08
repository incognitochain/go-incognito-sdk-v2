package incclient

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"math/big"
	"testing"
	"time"
)

//UTILS
const (
	pTokenID = "c7545459764224a000a9b323850648acf271186238210ce474b505cd17cc93a0"
	pEthID   = "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854"
)

var ic *IncClient

func initClients() error {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		return fmt.Errorf("cannot init new incognito client")
	}

	return nil
}

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

func DecodeBurnProof(r *jsonresult.InstructionProof) (*BurnProof, error) {
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
	for _, sIdx := range r.BeaconSigIndices {
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
	for _, sIdx := range r.BridgeSigIndices {
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

//END UTILS

//TEST FUNCTIONS
func TestIncClient_ShieldETH(t *testing.T) {
	err := initClients()
	if err != nil {
		panic(err)
	}

	//Incognito keys
	privateKey := ""

	oldBalance, err := ic.GetBalance(privateKey, pEthID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current balance of token %v: %v\n", pEthID, oldBalance)

	ethTxHash := "" // input the ETH transaction hash here.
	fmt.Printf("Start shielding eth...")

	ethProof, pETHAmount, err := ic.GetEVMDepositProof(ethTxHash)
	if err != nil {
		panic(err)
	}

	txHashStr, err := ic.CreateAndSendIssuingEVMRequestTransaction(privateKey, pEthID, *ethProof)
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
		newBalance, err := ic.GetBalance(privateKey, pEthID)
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
	err := initClients()
	if err != nil {
		panic(err)
	}

	//Incognito keys
	privateKey := ""

	oldBalance, err := ic.GetBalance(privateKey, pTokenID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current balance of token %v: %v\n", pTokenID, oldBalance)

	fmt.Printf("Start shielding token...\n")

	ethTxHash := "" // an ETH transaction
	ethProof, _, err := ic.GetEVMDepositProof(ethTxHash)
	if err != nil {
		panic(err)
	}

	txHashStr, err := ic.CreateAndSendIssuingEVMRequestTransaction(privateKey, pTokenID, *ethProof)
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
}

func TestIncClient_UnShieldETH(t *testing.T) {
	err := initClients()
	if err != nil {
		panic(err)
	}

	privateKey := ""
	remoteAddr := "" // an ETH address
	burnedAmount := uint64(50000000)
	burnedTxHash, err := ic.CreateAndSendBurningRequestTransaction(privateKey, remoteAddr, pEthID, burnedAmount)
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
	burningProofResult, err := ic.GetBurnProof(burnedTxHash)
	if err != nil {
		panic(err)
	}

	burnProof, err := DecodeBurnProof(burningProofResult)
	if err != nil {
		panic(err)
	}

	jsb, _ := json.Marshal(burnProof)

	fmt.Printf("Burn proof from Incog: %v\n", string(jsb))

	fmt.Printf("Finish getting the burning proof\n")
}

func TestIncClient_UnShieldERC20(t *testing.T) {
	err := initClients()
	if err != nil {
		panic(err)
	}

	privateKey := ""
	tokenIDStr := "c7545459764224a000a9b323850648acf271186238210ce474b505cd17cc93a0" //incognito tokenID for pDAI
	remoteAddr := ""                                                                 // an ETH address
	burnedAmount := uint64(10000000)
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
	burningProofResult, err := ic.GetBurnProof(burnedTxHash)
	if err != nil {
		panic(err)
	}

	burnProof, err := DecodeBurnProof(burningProofResult)
	if err != nil {
		panic(err)
	}

	jsb, _ := json.Marshal(burnProof)

	fmt.Printf("Burn proof from Incog: %v\n", string(jsb))
	fmt.Printf("Finish getting the burning proof\n")
}

func TestIncClient_GetETHTxReceipt(t *testing.T) {
	err := initClients()
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
	err := initClients()
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
