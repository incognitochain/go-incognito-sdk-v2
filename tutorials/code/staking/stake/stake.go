package main

import (
	"fmt"
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	privateSeed := incclient.PrivateKeyToMiningKey(privateKey) //NOTE: the private seed (a.k.a the mining key) can be randomly generated and not be dependent on the private key
	candidateAddress := incclient.PrivateKeyToPaymentAddress(privateKey, -1)
	rewardAddress := candidateAddress //NOTE: the reward receiver can either be the same as the candidate address or be different
	autoReStake := true
	beaconStakingAmount := uint64(1750 * 1e9 * 50)
	beaconAddStakingAmount := uint64(1750 * 1e9 * 5)

	txHash1, err := client.CreateAndSendShardStakingTransaction(privateKey, privateSeed, candidateAddress, rewardAddress, autoReStake)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("txHash shard staking %v\n", txHash1)

	txHash2, err := client.CreateAndSendBeaconStakingTransaction(privateKey, privateSeed, candidateAddress, rewardAddress, beaconStakingAmount)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("txHash beacon staking %v\n", txHash2)

	txHash3, err := client.CreateAndSendBeaconAddStakingTransaction(privateKey, privateSeed, candidateAddress, beaconAddStakingAmount)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("txHash beacon add staking %v\n", txHash3)
}
