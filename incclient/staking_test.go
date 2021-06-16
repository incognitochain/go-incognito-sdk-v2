package incclient

import (
	"fmt"
	"testing"
)

func TestIncClient_CreateAndSendShardStakingTransaction(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
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

func TestIncClient_CreateAndSendUnStakingTransaction(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
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

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	rewardAddress := PrivateKeyToPaymentAddress(privateKey, -1)

	txHash, err := ic.CreateAndSendWithDrawRewardTransaction(privateKey, rewardAddress)
	if err != nil {
		panic(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}
