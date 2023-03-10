package incclient

import (
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"github.com/pkg/errors"
)

const (
	SHARD_STAKING_AMOUNT      = DefaultShardStakeAmount
	MIN_BEACON_STAKING_AMOUNT = DefaultBeaconStakeAmount
)

// CreateShardStakingTransaction creates a raw staking transaction.
func (client *IncClient) CreateShardStakingTransaction(privateKey, privateSeed, candidateAddr, rewardReceiverAddr string, autoStake bool) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	funderAddr := senderWallet.Base58CheckSerialize(wallet.PaymentAddressType)

	if len(candidateAddr) == 0 {
		candidateAddr = funderAddr
	}
	if len(rewardReceiverAddr) == 0 {
		rewardReceiverAddr = funderAddr
	}

	candidateWallet, err := wallet.Base58CheckDeserialize(candidateAddr)
	if err != nil {
		return nil, "", err
	}
	pk := candidateWallet.KeySet.PaymentAddress.Pk
	if len(pk) == 0 {
		return nil, "", fmt.Errorf("candidate payment address invalid: %v", candidateAddr)
	}

	seed, _, err := base58.Base58Check{}.Decode(privateSeed)
	if err != nil {
		return nil, "", fmt.Errorf("cannot decode private seed: %v", privateSeed)
	}

	committeePK, err := key.NewCommitteeKeyFromSeed(seed, pk)
	if err != nil {
		return nil, "", fmt.Errorf("cannot create committee key from pk: %v, seed: %v. Error: %v", pk, seed, err)
	}

	committeePKBytes, err := committeePK.Bytes()
	if err != nil {
		return nil, "", fmt.Errorf("committee to bytes error: %v", err)
	}

	stakingAmount := SHARD_STAKING_AMOUNT

	stakingMetadata, err := metadata.NewStakingMetadata(metadata.ShardStakingMeta, funderAddr, rewardReceiverAddr, stakingAmount,
		base58.Base58Check{}.Encode(committeePKBytes, common.ZeroByte), autoStake)

	txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{stakingAmount}, 0, nil, stakingMetadata, nil)

	return client.CreateRawTransaction(txParam, -1)
}

func (client *IncClient) CreateBeaconStakingTransaction(privateKey, privateSeed, candidateAddr, rewardReceiverAddr string, stakingAmount uint64) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	funderAddr := senderWallet.Base58CheckSerialize(wallet.PaymentAddressType)

	if len(candidateAddr) == 0 {
		candidateAddr = funderAddr
	}
	if len(rewardReceiverAddr) == 0 {
		rewardReceiverAddr = funderAddr
	}

	candidateWallet, err := wallet.Base58CheckDeserialize(candidateAddr)
	if err != nil {
		return nil, "", err
	}
	pk := candidateWallet.KeySet.PaymentAddress.Pk
	if len(pk) == 0 {
		return nil, "", fmt.Errorf("candidate payment address invalid: %v", candidateAddr)
	}

	seed, _, err := base58.Base58Check{}.Decode(privateSeed)
	if err != nil {
		return nil, "", fmt.Errorf("cannot decode private seed: %v", privateSeed)
	}

	committeePK, err := key.NewCommitteeKeyFromSeed(seed, pk)
	if err != nil {
		return nil, "", fmt.Errorf("cannot create committee key from pk: %v, seed: %v. Error: %v", pk, seed, err)
	}

	committeePKBytes, err := committeePK.Bytes()
	if err != nil {
		return nil, "", fmt.Errorf("committee to bytes error: %v", err)
	}

	if (stakingAmount < MIN_BEACON_STAKING_AMOUNT) || (stakingAmount%SHARD_STAKING_AMOUNT != 0) {
		return nil, "", fmt.Errorf("Invalid beacon staking amount: %v, min beacon staking: %v, shard staking amount %v", stakingAmount, MIN_BEACON_STAKING_AMOUNT, SHARD_STAKING_AMOUNT)
	}

	stakingMetadata, err := metadata.NewStakingMetadata(metadata.BeaconStakingMeta, funderAddr, rewardReceiverAddr, stakingAmount,
		base58.Base58Check{}.Encode(committeePKBytes, common.ZeroByte), true)

	txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{stakingAmount}, 0, nil, stakingMetadata, nil)

	return client.CreateRawTransaction(txParam, -1)
}

func (client *IncClient) CreateBeaconAddStakingTransaction(privateKey, privateSeed, candidateAddr string, addStakingAmount uint64) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	funderAddr := senderWallet.Base58CheckSerialize(wallet.PaymentAddressType)

	if len(candidateAddr) == 0 {
		candidateAddr = funderAddr
	}

	candidateWallet, err := wallet.Base58CheckDeserialize(candidateAddr)
	if err != nil {
		return nil, "", err
	}
	pk := candidateWallet.KeySet.PaymentAddress.Pk
	if len(pk) == 0 {
		return nil, "", fmt.Errorf("candidate payment address invalid: %v", candidateAddr)
	}

	seed, _, err := base58.Base58Check{}.Decode(privateSeed)
	if err != nil {
		return nil, "", fmt.Errorf("cannot decode private seed: %v", privateSeed)
	}

	committeePK, err := key.NewCommitteeKeyFromSeed(seed, pk)
	if err != nil {
		return nil, "", fmt.Errorf("cannot create committee key from pk: %v, seed: %v. Error: %v", pk, seed, err)
	}
	committeePKStr, err := committeePK.ToBase58()
	if err != nil {
		return nil, "", fmt.Errorf("cannot create committee key from pk: %v, seed: %v. Error: %v", pk, seed, err)
	}

	if (addStakingAmount < SHARD_STAKING_AMOUNT*3) || (addStakingAmount%SHARD_STAKING_AMOUNT != 0) {
		return nil, "", fmt.Errorf("Invalid beacon staking amount: %v, min add staking amount: %v, shard staking amount %v", addStakingAmount, SHARD_STAKING_AMOUNT*3, SHARD_STAKING_AMOUNT)
	}

	addStakingMetadata, err := metadata.NewAddStakingMetadata(committeePKStr, addStakingAmount)

	txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{addStakingAmount}, 0, nil, addStakingMetadata, nil)

	return client.CreateRawTransaction(txParam, -1)
}

// CreateAndSendShardStakingTransaction creates a raw staking transaction and broadcasts it to the blockchain.
func (client *IncClient) CreateAndSendShardStakingTransaction(privateKey, privateSeed, candidateAddr, rewardReceiverAddr string, autoStake bool) (string, error) {
	encodedTx, txHash, err := client.CreateShardStakingTransaction(privateKey, privateSeed, candidateAddr, rewardReceiverAddr, autoStake)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

func (client *IncClient) CreateAndSendBeaconStakingTransaction(privateKey, privateSeed, candidateAddr, rewardReceiverAddr string, stakingAmount uint64) (string, error) {
	encodedTx, txHash, err := client.CreateBeaconStakingTransaction(privateKey, privateSeed, candidateAddr, rewardReceiverAddr, stakingAmount)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

func (client *IncClient) CreateAndSendBeaconAddStakingTransaction(privateKey, privateSeed, candidateAddr string, stakingAmount uint64) (string, error) {
	encodedTx, txHash, err := client.CreateBeaconAddStakingTransaction(privateKey, privateSeed, candidateAddr, stakingAmount)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreateUnStakingTransaction creates a raw un-staking transaction.
func (client *IncClient) CreateUnStakingTransaction(privateKey, privateSeed, candidateAddr string) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	funderAddr := senderWallet.Base58CheckSerialize(wallet.PaymentAddressType)

	if len(candidateAddr) == 0 {
		candidateAddr = funderAddr
	}

	candidateWallet, err := wallet.Base58CheckDeserialize(candidateAddr)
	if err != nil {
		return nil, "", err
	}
	pk := candidateWallet.KeySet.PaymentAddress.Pk
	if len(pk) == 0 {
		return nil, "", fmt.Errorf("candidate payment address invalid: %v", candidateAddr)
	}

	seed, _, err := base58.Base58Check{}.Decode(privateSeed)
	if err != nil {
		return nil, "", fmt.Errorf("cannot decode private seed %v: %v", privateSeed, err)
	}

	committeePK, err := key.NewCommitteeKeyFromSeed(seed, pk)
	if err != nil {
		return nil, "", fmt.Errorf("cannot create committee key from pk: %v, seed: %v. Error: %v", pk, seed, err)
	}

	committeePKBytes, err := committeePK.Bytes()
	if err != nil {
		return nil, "", fmt.Errorf("committee to bytes error: %v", err)
	}
	unStakingMetadata, err := metadata.NewUnStakingMetadata(base58.Base58Check{}.Encode(committeePKBytes, common.ZeroByte))

	txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{0}, 0, nil, unStakingMetadata, nil)

	return client.CreateRawTransaction(txParam, -1)
}

// CreateAndSendUnStakingTransaction creates a raw un-staking transaction and broadcasts it to the blockchain.
func (client *IncClient) CreateAndSendUnStakingTransaction(privateKey, privateSeed, candidateAddr string) (string, error) {
	encodedTx, txHash, err := client.CreateUnStakingTransaction(privateKey, privateSeed, candidateAddr)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreateWithDrawRewardTransaction creates a raw reward-withdrawing transaction.
func (client *IncClient) CreateWithDrawRewardTransaction(privateKey, addr, tokenIDStr string, version int8) ([]byte, string, error) {
	if version != 1 && version != 2 {
		return nil, "", fmt.Errorf("only support version 1 or 2")
	}

	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		Logger.Printf("%v\n", err)
		return nil, "", err
	}

	funderAddr := senderWallet.Base58CheckSerialize(wallet.PaymentAddressType)
	if len(addr) == 0 {
		addr = funderAddr
	}
	if version == 1 {
		addr, err = wallet.GetPaymentAddressV1(addr, false)
		if err != nil {
			Logger.Printf("%v\n", err)
			return nil, "", err
		}
	}

	if len(tokenIDStr) == 0 {
		Logger.Printf("No tokenID provided, using the default PRV\n")
		tokenIDStr = common.PRVIDStr
	}

	withdrawRewardMetadata, err := metadata.NewWithDrawRewardRequest(tokenIDStr, addr, 0, metadata.WithDrawRewardRequestMeta)

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, nil, withdrawRewardMetadata, nil)

	return client.CreateRawTransaction(txParam, version)
}

// CreateAndSendWithDrawRewardTransaction creates a raw reward-withdrawing transaction and broadcasts it to the blockchain.
func (client *IncClient) CreateAndSendWithDrawRewardTransaction(privateKey, addr, tokenIDStr string, version int8) (string, error) {
	encodedTx, txHash, err := client.CreateWithDrawRewardTransaction(privateKey, addr, tokenIDStr, version)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

func (client *IncClient) CreateReDelegateTransaction(privateKey, privateSeed, candidateAddr string, delegate string) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	funderAddr := senderWallet.Base58CheckSerialize(wallet.PaymentAddressType)

	if len(candidateAddr) == 0 {
		candidateAddr = funderAddr
	}

	candidateWallet, err := wallet.Base58CheckDeserialize(candidateAddr)
	if err != nil {
		return nil, "", err
	}
	pk := candidateWallet.KeySet.PaymentAddress.Pk
	if len(pk) == 0 {
		return nil, "", fmt.Errorf("candidate payment address invalid: %v", candidateAddr)
	}

	seed, _, err := base58.Base58Check{}.Decode(privateSeed)
	if err != nil {
		return nil, "", fmt.Errorf("cannot decode private seed: %v", privateSeed)
	}

	committeePK, err := key.NewCommitteeKeyFromSeed(seed, pk)
	if err != nil {
		return nil, "", fmt.Errorf("cannot create committee key from pk: %v, seed: %v. Error: %v", pk, seed, err)
	}
	committeePKStr, err := committeePK.ToBase58()
	if err != nil {
		return nil, "", fmt.Errorf("cannot create committee key from pk: %v, seed: %v. Error: %v", pk, seed, err)
	}

	delegatePKStruct := &key.CommitteePublicKey{}
	err = delegatePKStruct.FromString(delegate)
	if err != nil {
		return nil, "", err
	}
	delegateUIDI, err := client.rpcServer.GetBeaconCandidateUID(delegate)
	if err != nil {
		return nil, "", err
	}

	delegateUID, ok := delegateUIDI.(string)
	if !ok {
		return nil, "", errors.Errorf("Expected get Beacon Candidate UID of beacon %+v at string, but received %+v", delegate, delegateUIDI)
	}
	redelegateMetadata, err := metadata.NewReDelegateMetadata(committeePKStr, delegate, delegateUID)

	txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{0}, 0, nil, redelegateMetadata, nil)

	return client.CreateRawTransaction(txParam, -1)
}

func (client *IncClient) CreateAndSendReDelegateTransaction(privateKey, privateSeed, candidateAddr string, delegate string) (string, error) {
	encodedTx, txHash, err := client.CreateReDelegateTransaction(privateKey, privateSeed, candidateAddr, delegate)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}
