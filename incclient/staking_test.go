package incclient

import (
	"fmt"
	"testing"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

func TestIncClient_CreateAndSendShardStakingTransaction(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	privateKey := ""
	privateSeed := PrivateKeyToMiningKey(privateKey) //NOTE: the private seed (a.k.a the mining key) can be randomly generated and not be dependent on the private key
	candidateAddress := PrivateKeyToPaymentAddress(privateKey, -1)
	rewardAddress := candidateAddress //NOTE: the reward receiver can either be the same as the candidate address or be different
	autoReStake := true

	txHash, err := ic.CreateAndSendShardStakingTransaction(privateKey, privateSeed, candidateAddress, rewardAddress, autoReStake)
	if err != nil {
		panic(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}

func TestIncClient_CreateAndSendBeaconStakingTransaction(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	privateKey := ""
	privateSeed := PrivateKeyToMiningKey(privateKey) //NOTE: the private seed (a.k.a the mining key) can be randomly generated and not be dependent on the private key
	candidateAddress := PrivateKeyToPaymentAddress(privateKey, -1)
	rewardAddress := candidateAddress //NOTE: the reward receiver can either be the same as the candidate address or be different
	stakingAmount := uint64(1750 * 1e9 * 70)

	txHash, err := ic.CreateAndSendBeaconStakingTransaction(privateKey, privateSeed, candidateAddress, rewardAddress, stakingAmount)
	if err != nil {
		panic(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}

func TestIncClient_CreateAndSendBeaconAddStakingTransaction(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	privateKey := ""
	privateSeed := PrivateKeyToMiningKey(privateKey) //NOTE: the private seed (a.k.a the mining key) can be randomly generated and not be dependent on the private key
	candidateAddress := PrivateKeyToPaymentAddress(privateKey, -1)
	addStakingAmount := uint64(1750 * 1e9 * 7)

	txHash, err := ic.CreateAndSendBeaconAddStakingTransaction(privateKey, privateSeed, candidateAddress, addStakingAmount)
	if err != nil {
		panic(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}

func TestIncClient_CreateAndSendUnStakingTransaction(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	privateKey := ""
	privateSeed := PrivateKeyToMiningKey(privateKey) //NOTE: the private seed (a.k.a the mining key) can be randomly generated and not be dependent on the private key
	candidateAddress := PrivateKeyToPaymentAddress(privateKey, -1)

	txHash, err := ic.CreateAndSendUnStakingTransaction(privateKey, privateSeed, candidateAddress)
	if err != nil {
		panic(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}

func TestIncClient_CreateAndSendWithDrawRewardTransaction(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	privateKey := ""
	rewardAddress := PrivateKeyToPaymentAddress(privateKey, -1)

	txHash, err := ic.CreateAndSendWithDrawRewardTransaction(privateKey, rewardAddress, common.PRVIDStr, 2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}
