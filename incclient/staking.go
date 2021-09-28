package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
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

	stakingAmount := uint64(1750000000000)

	stakingMetadata, err := metadata.NewStakingMetadata(metadata.ShardStakingMeta, funderAddr, rewardReceiverAddr, stakingAmount,
		base58.Base58Check{}.Encode(committeePKBytes, common.ZeroByte), autoStake)

	txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{stakingAmount}, 0, nil, stakingMetadata, nil)

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

	unStakingMetadata, err := metadata.NewStopAutoStakingMetadata(metadata.StopAutoStakingMeta, base58.Base58Check{}.Encode(committeePKBytes, common.ZeroByte))

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
