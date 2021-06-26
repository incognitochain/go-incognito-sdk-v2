package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
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

	txHash, err := client.CreateAndSendShardStakingTransaction(privateKey, privateSeed, candidateAddress, rewardAddress, autoReStake)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}
