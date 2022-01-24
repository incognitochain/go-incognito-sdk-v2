package rpc

// rpc cmd method
const (
	// test rpc server
	testHttpServer = "testrpcserver"
	startProfiling = "startprofiling"
	stopProfiling  = "stopprofiling"
	exportMetrics  = "exportmetrics"

	getNetworkInfo       = "getnetworkinfo"
	getConnectionCount   = "getconnectioncount"
	getAllConnectedPeers = "getallconnectedpeers"
	getAllPeers          = "getallpeers"
	getNodeRole          = "getnoderole"
	getInOutMessages     = "getinoutmessages"
	getInOutMessageCount = "getinoutmessagecount"

	estimateFee              = "estimatefee"
	estimateFeeV2            = "estimatefeev2"
	estimateFeeWithEstimator = "estimatefeewithestimator"

	getActiveShards    = "getactiveshards"
	getMaxShardsNumber = "getmaxshardsnumber"

	getMiningInfo                 = "getmininginfo"
	getSyncStats                  = "getsyncstats"
	getRawMempool                 = "getrawmempool"
	getNumberOfTxsInMempool       = "getnumberoftxsinmempool"
	getMempoolEntry               = "getmempoolentry"
	removeTxInMempool             = "removetxinmempool"
	getBeaconPoolState            = "getbeaconpoolstate"
	getShardPoolState             = "getshardpoolstate"
	getShardPoolLatestValidHeight = "getshardpoollatestvalidheight"
	//getShardToBeaconPoolState     = "getshardtobeaconpoolstate"
	//getCrossShardPoolState        = "getcrossshardpoolstate"
	getNextCrossShard           = "getnextcrossshard"
	getShardToBeaconPoolStateV2 = "getshardtobeaconpoolstatev2"
	getCrossShardPoolStateV2    = "getcrossshardpoolstatev2"
	getShardPoolStateV2         = "getshardpoolstatev2"
	getBeaconPoolStateV2        = "getbeaconpoolstatev2"
	//getFeeEstimator             = "getfeeestimator"
	setBackup                   = "setbackup"
	getLatestBackup             = "getlatestbackup"
	getBestBlock                = "getbestblock"
	getBestBlockHash            = "getbestblockhash"
	getBlocks                   = "getblocks"
	retrieveBlock               = "retrieveblock"
	retrieveBlockByHeight       = "retrieveblockbyheight"
	retrieveBeaconBlock         = "retrievebeaconblock"
	retrieveBeaconBlockByHeight = "retrievebeaconblockbyheight"
	getBlockChainInfo           = "getblockchaininfo"
	getBlockCount               = "getblockcount"
	getBlockHash                = "getblockhash"

	listOutputCoins                            = "listoutputcoins"
	listOutputCoinsFromCache                   = "listoutputcoinsfromcache"
	listOutputTokens                           = "listoutputtokens"
	createRawTransaction                       = "createtransaction"
	sendRawTransaction                         = "sendtransaction"
	createAndSendTransaction                   = "createandsendtransaction"
	createandsendconversiontransaction         = "createandsendconversiontransaction"
	createAndSendCustomTokenTransaction        = "createandsendcustomtokentransaction"
	sendRawCustomTokenTransaction              = "sendrawcustomtokentransaction"
	createRawCustomTokenTransaction            = "createrawcustomtokentransaction"
	createConvertCoinVer1ToVer2TxToken         = "createconvertcoinver1tover2txtoken"
	createRawPrivacyCustomTokenTransaction     = "createrawprivacycustomtokentransaction"
	sendRawPrivacyCustomTokenTransaction       = "sendrawprivacycustomtokentransaction"
	createAndSendPrivacyCustomTokenTransaction = "createandsendprivacycustomtokentransaction"
	getMempoolInfo                             = "getmempoolinfo"
	getPendingTxsInBlockgen                    = "getpendingtxsinblockgen"
	getCandidateList                           = "getcandidatelist"
	getCommitteeList                           = "getcommitteelist"
	canPubkeyStake                             = "canpubkeystake"
	getTotalTransaction                        = "gettotaltransaction"
	listUnspentCustomToken                     = "listunspentcustomtoken"
	getBalanceCustomToken                      = "getbalancecustomtoken"
	getTransactionByHash                       = "gettransactionbyhash"
	getEncodedTransactionsByHashes             = "getencodedtransactionsbyhashes"
	gettransactionhashbyreceiver               = "gettransactionhashbyreceiver"
	gettransactionhashbyreceiverv2             = "gettransactionhashbyreceiverv2"
	gettransactionbyreceiver                   = "gettransactionbyreceiver"
	gettransactionbyreceiverv2                 = "gettransactionbyreceiverv2"
	gettransactionbyserialnumber               = "gettransactionbyserialnumber"
	gettransactionbypublickey                  = "gettransactionbypublickey"
	listCustomToken                            = "listcustomtoken"
	listPrivacyCustomToken                     = "listprivacycustomtoken"
	listPrivacyCustomTokenIDs                  = "listprivacycustomtokenids"
	getPrivacyCustomToken                      = "getprivacycustomtoken"
	listPrivacyCustomTokenByShard              = "listprivacycustomtokenbyshard"
	getBalancePrivacyCustomToken               = "getbalanceprivacycustomtoken"
	listUnspentOutputTokens                    = "listunspentoutputtokens"
	customTokenTxs                             = "customtoken"
	listCustomTokenHolders                     = "customtokenholder"
	privacyCustomTokenTxs                      = "privacycustomtoken"
	checkHashValue                             = "checkhashvalue"
	getListCustomTokenBalance                  = "getlistcustomtokenbalance"
	getListPrivacyCustomTokenBalance           = "getlistprivacycustomtokenbalance"
	getBlockHeader                             = "getheader"
	getCrossShardBlock                         = "getcrossshardblock"
	randomCommitments                          = "randomcommitments"
	hasSerialNumbers                           = "hasserialnumbers"
	hasSerialNumbersInMempool                  = "hasserialnumbersinmempool"
	hasSnDerivators                            = "hassnderivators"
	listSnDerivators                           = "listsnderivators"
	listSerialNumbers                          = "listserialnumbers"
	listCommitments                            = "listcommitments"
	listCommitmentIndices                      = "listcommitmentindices"
	createAndSendStakingTransaction            = "createandsendstakingtransaction"
	createAndSendStopAutoStakingTransaction    = "createandsendstopautostakingtransaction"
	createAndSendTokenInitTransaction          = "createandsendtokeninittransaction"
	decryptoutputcoinbykeyoftransaction        = "decryptoutputcoinbykeyoftransaction"
	randomCommitmentsAndPublicKeys             = "randomcommitmentsandpublickeys"
	getOTACoinLength                           = "getotacoinlength"
	getOTACoinsByIndices                       = "getotacoinsbyindices"

	//===========For Testing and Benchmark==============
	getAndSendTxsFromFile      = "getandsendtxsfromfile"
	getAndSendTxsFromFileV2    = "getandsendtxsfromfilev2"
	unlockMempool              = "unlockmempool"
	handleGetConsensusInfoV3   = "getconsensusinfov3"
	getAutoStakingByHeight     = "getautostakingbyheight"
	getCommitteeState          = "getcommitteestate"
	convertPaymentAddress      = "convertpaymentaddress"
	getCommitteeStateByShard   = "getcommitteestatebyshard"
	getSlashingCommittee       = "getslashingcommittee"
	getSlashingCommitteeDetail = "getslashingcommitteedetail"
	getRewardAmountByEpoch     = "getrewardamountbyepoch"
	//==================================================

	getShardBestState        = "getshardbeststate"
	getShardBestStateDetail  = "getshardbeststatedetail"
	getBeaconBestState       = "getbeaconbeststate"
	getBeaconBestStateDetail = "getbeaconbeststatedetail"

	// Wallet rpc cmd
	listAccounts                    = "listaccounts"
	getAccount                      = "getaccount"
	getAddressesByAccount           = "getaddressesbyaccount"
	getAccountAddress               = "getaccountaddress"
	dumpPrivkey                     = "dumpprivkey"
	importAccount                   = "importaccount"
	removeAccount                   = "removeaccount"
	listUnspentOutputCoins          = "listunspentoutputcoins"
	listUnspentOutputCoinsFromCache = "listunspentoutputcoinsfromcache"
	getBalance                      = "getbalance"
	getBalanceByPrivatekey          = "getbalancebyprivatekey"
	getBalanceByPaymentAddress      = "getbalancebypaymentaddress"
	getReceivedByAccount            = "getreceivedbyaccount"
	setTxFee                        = "settxfee"
	submitKey                       = "submitkey"
	authorizedSubmitKey             = "authorizedsubmitkey"
	getKeySubmissionInfo            = "getkeysubmissioninfo"

	// walletsta
	getPublicKeyFromPaymentAddress = "getpublickeyfrompaymentaddress"
	defragmentAccount              = "defragmentaccount"
	defragmentAccountV2            = "defragmentaccountv2"
	defragmentAccountToken         = "defragmentaccounttoken"
	defragmentAccountTokenV2       = "defragmentaccounttokenv2"

	getStackingAmount = "getstackingamount"

	// utils
	hashToIdenticon = "hashtoidenticon"
	generateTokenID = "generatetokenid"

	createIssuingRequest               = "createissuingrequest"
	sendIssuingRequest                 = "sendissuingrequest"
	createAndSendIssuingRequest        = "createandsendissuingrequest"
	createAndSendIssuingRequestV2      = "createandsendissuingrequestv2"
	createAndSendContractingRequest    = "createandsendcontractingrequest"
	createAndSendContractingRequestV2  = "createandsendcontractingrequestv2"
	createAndSendBurningRequest        = "createandsendburningrequest"
	createAndSendBurningRequestV2      = "createandsendburningrequestv2"
	createAndSendTxWithIssuingETHReq   = "createandsendtxwithissuingethreq"
	createAndSendTxWithIssuingETHReqV2 = "createandsendtxwithissuingethreqv2"
	checkETHHashIssued                 = "checkethhashissued"
	getAllBridgeTokens                 = "getallbridgetokens"
	getETHHeaderByHash                 = "getethheaderbyhash"
	getBridgeReqWithStatus             = "getbridgereqwithstatus"

	// Incognito -> Ethereum bridge
	getBeaconSwapProof       = "getbeaconswapproof"
	getLatestBeaconSwapProof = "getlatestbeaconswapproof"
	getBridgeSwapProof       = "getbridgeswapproof"
	getLatestBridgeSwapProof = "getlatestbridgeswapproof"
	getBurnProof             = "getburnproof"
	getBSCBurnProof          = "getbscburnproof"
	getPRVERC20BurnProof     = "getprverc20burnproof"
	getPRVBEP20BurnProof     = "getprvbep20burnproof"

	// reward
	CreateRawWithDrawTransaction = "withdrawreward"
	getRewardAmount              = "getrewardamount"
	getRewardAmountByPublicKey   = "getrewardamountbypublickey"
	listRewardAmount             = "listrewardamount"

	revertbeaconchain = "revertbeaconchain"
	revertshardchain  = "revertshardchain"

	enableMining                = "enablemining"
	getChainMiningStatus        = "getchainminingstatus"
	getPublickeyMining          = "getpublickeymining"
	getPublicKeyRole            = "getpublickeyrole"
	getRoleByValidatorKey       = "getrolebyvalidatorkey"
	getIncognitoPublicKeyRole   = "getincognitopublickeyrole"
	getMinerRewardFromMiningKey = "getminerrewardfromminingkey"

	// slash
	getProducersBlackList       = "getproducersblacklist"
	getProducersBlackListDetail = "getproducersblacklistdetail"

	// pde
	getPDEState                                = "getpdestate"
	createAndSendTxWithWithdrawalReq           = "createandsendtxwithwithdrawalreq"
	createAndSendTxWithWithdrawalReqV2         = "createandsendtxwithwithdrawalreqv2"
	createAndSendTxWithPDEFeeWithdrawalReq     = "createandsendtxwithpdefeewithdrawalreq"
	createAndSendTxWithPTokenTradeReq          = "createandsendtxwithptokentradereq"
	createAndSendTxWithPTokenCrossPoolTradeReq = "createandsendtxwithptokencrosspooltradereq"
	createAndSendTxWithPRVTradeReq             = "createandsendtxwithprvtradereq"
	createAndSendTxWithPRVCrossPoolTradeReq    = "createandsendtxwithprvcrosspooltradereq"
	createAndSendTxWithPTokenContribution      = "createandsendtxwithptokencontribution"
	createAndSendTxWithPRVContribution         = "createandsendtxwithprvcontribution"
	createAndSendTxWithPTokenContributionV2    = "createandsendtxwithptokencontributionv2"
	createAndSendTxWithPRVContributionV2       = "createandsendtxwithprvcontributionv2"
	convertNativeTokenToPrivacyToken           = "convertnativetokentoprivacytoken"
	convertPrivacyTokenToNativeToken           = "convertprivacytokentonativetoken"
	getPDEContributionStatus                   = "getpdecontributionstatus"
	getPDEContributionStatusV2                 = "getpdecontributionstatusv2"
	getPDETradeStatus                          = "getpdetradestatus"
	getPDEWithdrawalStatus                     = "getpdewithdrawalstatus"
	getPDEFeeWithdrawalStatus                  = "getpdefeewithdrawalstatus"
	convertPDEPrices                           = "convertpdeprices"
	extractPDEInstsFromBeaconBlock             = "extractpdeinstsfrombeaconblock"

	pdexv3MintNft                                  = "pdexv3_txMintNft"
	getPdexv3State                                 = "pdexv3_getState"
	createAndSendTxWithPdexv3ModifyParams          = "pdexv3_txModifyParams"
	getPdexv3ParamsModifyingStatus                 = "pdexv3_getParamsModifyingStatus"
	pdexv3AddLiquidityV3                           = "pdexv3_txAddLiquidity"
	pdexv3WithdrawLiquidityV3                      = "pdexv3_txWithdrawLiquidity"
	getPdexv3ContributionStatus                    = "pdexv3_getContributionStatus"
	getPdexv3WithdrawLiquidityStatus               = "pdexv3_getWithdrawLiquidityStatus"
	getPdexv3MintNftStatus                         = "pdexv3_getMintNftStatus"
	pdexv3TxTrade                                  = "pdexv3_txTrade"
	pdexv3TxAddOrder                               = "pdexv3_txAddOrder"
	pdexv3TxWithdrawOrder                          = "pdexv3_txWithdrawOrder"
	pdexv3GetTradeStatus                           = "pdexv3_getTradeStatus"
	pdexv3GetAddOrderStatus                        = "pdexv3_getAddOrderStatus"
	pdexv3GetWithdrawOrderStatus                   = "pdexv3_getWithdrawOrderStatus"
	pdexv3Staking                                  = "pdexv3_txStake"
	pdexv3Unstaking                                = "pdexv3_txUnstake"
	pdexv3GetStakingStatus                         = "pdexv3_getStakingStatus"
	pdexv3GetUnstakingStatus                       = "pdexv3_getUnstakingStatus"
	getPdexv3EstimatedLPValue                      = "pdexv3_getEstimatedLPValue"
	getPdexv3EstimatedLPPoolReward                 = "pdexv3_getEstimatedLPPoolReward"
	createAndSendTxWithPdexv3WithdrawLPFee         = "pdexv3_txWithdrawLPFee"
	getPdexv3WithdrawalLPFeeStatus                 = "pdexv3_getWithdrawalLPFeeStatus"
	createAndSendTxWithPdexv3WithdrawProtocolFee   = "pdexv3_txWithdrawProtocolFee"
	getPdexv3WithdrawalProtocolFeeStatus           = "pdexv3_getWithdrawalProtocolFeeStatus"
	getPdexv3EstimatedStakingReward                = "pdexv3_getEstimatedStakingReward"
	getPdexv3EstimatedStakingPoolReward            = "pdexv3_getEstimatedStakingPoolReward"
	createAndSendTxWithPdexv3WithdrawStakingReward = "pdexv3_txWithdrawStakingReward"
	getPdexv3WithdrawalStakingRewardStatus         = "pdexv3_getWithdrawalStakingRewardStatus"

	// get burning address
	getBurningAddress = "getburningaddress"

	// portal
	getPortalV4State                           = "getportalv4state"
	getPortalV4Params                          = "getportalv4params"
	createAndSendTxWithShieldingRequest        = "createandsendtxshieldingrequest"
	getPortalShieldingRequestStatus            = "getportalshieldingrequeststatus"
	createAndSendTxWithPortalV4UnshieldRequest = "createandsendtxwithportalv4unshieldrequest"
	getPortalUnShieldingRequestStatus          = "getportalunshieldrequeststatus"
	getPortalBatchUnShieldingRequestStatus     = "getportalbatchunshieldrequeststatus"
	getSignedRawTransactionByBatchID           = "getportalsignedrawtransaction"
	createAndSendTxWithPortalReplacementFee    = "createandsendtxwithportalreplacebyfee"
	getPortalReplacementFeeStatus              = "getportalreplacebyfeestatus"
	createAndSendTxWithPortalSubmitConfirmedTx = "createandsendtxwithportalsubmitconfirmedtx"
	getPortalSubmitConfirmedTx                 = "getportalsubmitconfirmedtxstatus"
	getSignedRawReplaceFeeTransaction          = "getportalsignedrawreplacebyfeetransaction"
	createAndSendTxPortalConvertVaultRequest   = "createandsendtxportalconvertvault"
	getPortalConvertVaultTxStatus              = "getportalconvertvaultstatus"
	generatePortalShieldMultisigAddress        = "generateportalshieldmultisigaddress"
	generateOTDepositKey                       = "generateotdepositkey"
	getNextOTDepositKey                        = "getnextotdepositkey"
	hasOTDepositPubKey                         = "hasotdepositpubkey"
	getDepositTxsByPubKeys                     = "getdeposittxsbypubkeys"

	// relaying
	createAndSendTxWithRelayingBNBHeader = "createandsendtxwithrelayingbnbheader"
	createAndSendTxWithRelayingBTCHeader = "createandsendtxwithrelayingbtcheader"
	getRelayingBNBHeaderState            = "getrelayingbnbheaderstate"
	getRelayingBNBHeaderByBlockHeight    = "getrelayingbnbheaderbyblockheight"
	getBTCRelayingBestState              = "getbtcrelayingbeststate"
	getBTCBlockByHash                    = "getbtcblockbyhash"
	getLatestBNBHeaderBlockHeight        = "getlatestbnbheaderblockheight"

	// incognito mode for sc
	getBurnProofForDepositToSC                  = "getburnprooffordeposittosc"
	createAndSendBurningForDepositToSCRequest   = "createandsendburningfordeposittoscrequest"
	createAndSendBurningForDepositToSCRequestV2 = "createandsendburningfordeposittoscrequestv2"

	getBeaconPoolInfo     = "getbeaconpoolinfo"
	getShardPoolInfo      = "getshardpoolinfo"
	getCrossShardPoolInfo = "getcrossshardpoolinfo"
	getAllView            = "getallview"
	getAllViewDetail      = "getallviewdetail"

	// feature rewards
	getRewardFeature = "getrewardfeature"

	getTotalStaker = "gettotalstaker"

	//validator state
	getValKeyState = "getvalkeystate"

	getAllTradesInMemPool = "getalltradesinmempool"
	getAllTradesByAddress = "getalltradesbyaddress"
)
