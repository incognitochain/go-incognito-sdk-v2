package wallet

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/incognitochain/go-incognito-sdk-v2/key"

	"github.com/stretchr/testify/assert"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/tyler-smith/go-bip39"
)

var (
	numTests = 100
)

type vector struct {
	entropy      string
	seed         string
	mnemonic     string
	masterHexKey string
	password     string
	firstKey     string
	secondKey    string
}

func testMnemonicVectors() []vector { // test vectors built from https://iancoleman.io/bip39/#english
	return []vector{
		{
			entropy:  "00000000000000000000000000000000",
			mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
			seed:     "5eb00bbddcf069084889a8ab9155568165f5c453ccb85e70811aaed6f6da5fc19a5ac40b389cd370d086206dec8aa6c43daea6690f20ad3d8d48b2d2ce9e38e4",
		},
		{
			entropy:  "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
			mnemonic: "legal winner thank year wave sausage worth useful legal winner thank yellow",
			seed:     "878386efb78845b3355bd15ea4d39ef97d179cb712b77d5c12b6be415fffeffe5f377ba02bf3f8544ab800b955e51fbff09828f682052a20faa6addbbddfb096",
		},
		{
			entropy:  "9a8da4c1e5e8f874b9a2ae8afca812ad",
			mnemonic: "once honey corn slim moon demise track fiction memory torch again focus",
			seed:     "4735d20f5b5f4100fd1cea9bb9157481b9670bf98955bd4a7bf643df0e0f4cb59575ea477a040a7595aa41096bd3850b032d6960a198eac5b749d2ba5d0de827",
		},
		{
			entropy:  "ffffffffffffffffffffffffffffffff",
			mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo wrong",
			seed:     "b6a6d8921942dd9806607ebc2750416b289adea669198769f2e15ed926c3aa92bf88ece232317b4ea463e84b0fcd3b53577812ee449ccc448eb45e6f544e25b6",
		},
		{
			entropy:  "9894f9bb68dabc9d85cd4c9ff8c1a418",
			mnemonic: "obtain pond human spider profit excite blame praise paper ship harbor cradle",
			seed:     "c4e96a240d17b629dd9f8965cd6e6588a4abf6558173af31a05146687f93e2e647072f5f48154336fd684bd4df1da12402ab709bc745bd1e5362e2df577586ce",
		},
		{
			entropy:  "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
			mnemonic: "legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal will",
			seed:     "b059400ce0f55498a5527667e77048bb482ff6daa16c37b4b9e8af70c85b3f4df588004f19812a1a027c9a51e5e94259a560268e91cd10e206451a129826e740",
		},
		{
			entropy:  "808080808080808080808080808080808080808080808080",
			mnemonic: "letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter always",
			seed:     "04d5f77103510c41d610f7f5fb3f0badc77c377090815cee808ea5d2f264fdfabf7c7ded4be6d4c6d7cdb021ba4c777b0b7e57ca8aa6de15aeb9905dba674d66",
		},
		{
			entropy:  "ffffffffffffffffffffffffffffffffffffffffffffffff",
			mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo when",
			seed:     "d2911131a6dda23ac4441d1b66e2113ec6324354523acfa20899a2dcb3087849264e91f8ec5d75355f0f617be15369ffa13c3d18c8156b97cd2618ac693f759f",
		},
		{
			entropy:  "0000000000000000000000000000000000000000000000000000000000000000",
			mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
			seed:     "408b285c123836004f4b8842c89324c1f01382450c0d439af345ba7fc49acf705489c6fc77dbd4e3dc1dd8cc6bc9f043db8ada1e243c4a0eafb290d399480840",
		},
		{
			entropy:  "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
			mnemonic: "legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth title",
			seed:     "761914478ebf6fe16185749372e91549361af22b386de46322cf8b1ba7e92e80c4af05196f742be1e63aab603899842ddadf4e7248d8e43870a4b6ff9bf16324",
		},
		{
			entropy:  "31f4512bf0c1ab12e15e500c31de9b4b7ff3681764c041958660971872b10ae0",
			mnemonic: "cradle penalty enough thunder boy maximum lyrics skate around moment plug notice you reduce frown object dose promote oblige comfort mango flat clog cable",
			seed:     "20b2ad4777d04045bfecad79e5966dc1652191201a4c2ad093e1456fc93bd90f24968a4f3c28059dfdeb4f8f2a52ee4c3f14fb1bb8837638da3137190319de46",
		},
		{
			entropy:  "04a43d6973bf6e5d57711bb44574476b620859a1",
			mnemonic: "again capable fog trash want concert fruit casino reform close balcony strategy calm coast man",
			seed:     "2b94b9ccfb35bd27997024c9ca8d7ade65290285c0a7a4d3ec65ea19e24fcb448128272d5d77f49687a7ef668f290377b9daf1bac68b7a951252ebd6c609e21c",
		},
		{
			entropy:  "2784db83716fc2e9ab7f9e935d7dd95f5fbcc493",
			mnemonic: "chest chat this tissue wire inquiry pupil video nation typical iron salt wink girl estate",
			seed:     "fe4d79210b0effee1bab4912ab648f90eaadcb20a0147eab1c70c86522500b5a646b90283bb8a420689dc0fe2260f223307d233616fcf3f6a59c581793f1fe8a",
		},
		{
			entropy:  "017773efff8e311802b6f84860a2d321b4b42988c1944ab08b04e9c1",
			mnemonic: "accident romance winner yellow toast metal approve tenant embody agree regular drop ensure below cart crater enjoy lounge scorpion squeeze beef",
			seed:     "0d2a122edcdd236855df9cd1db587b4e80578955476498fd7098e6b713a992117a6326fa4a6d4bc39a5ffff27d5dc1c540d46c6df99bac090e6bd2f4e1a23035",
		},
		{
			entropy:  "f4931e3d48fe05307f9f286aa7932405711841076a3376d6c7f45077",
			mnemonic: "virus ocean monster music theory oblige wrist topic height develop simple april bag calm buffalo edit item renew wonder path imitate",
			seed:     "067108e9bcd388fcfb35c0975c2f5a06da8058dc35a9de4059e0cd0a4cd9e3c784c0117b28b03b6ee2704d38dfeefd9cd0c4e9958d4ca9e45a39d8971f6d7570",
		},
		{
			entropy:  "e8f66170ac714e58637d1657470d10adf80f6f68a7f088f3bb5c4fcd",
			mnemonic: "trophy reason found flight belt club mistake people firm debris during fossil liberty response penalty wrap material oven strike panel stuff",
			seed:     "212de7dedcea264a02743dfd61d12e35c31ae4f252d701bf066d6bb615bd9a10457ceb5285388258da6133e82ec89141e67f961e29016670189d59e5cc6427bf",
		},
		{
			entropy:  "401ea7721f228febef8d60489ae202aa4df7bb610f3a525926fea69e",
			mnemonic: "divorce vivid symptom dinner cigar vote sail project embrace strike level fee term tank loud trap false since sausage essay tide",
			seed:     "f16d007b95303418e709b5985c3dc97455772d4789576929756146f8c11dbc44379acbd9df717413015709d928f1c0288c24e1ad1dcb3c8e0f1961aee014bf05",
		},
		{
			entropy:  "8ce34a72726e62a3243f6c85b1b41b69",
			mnemonic: "mind bottom orient tooth tower face movie unique mad misery almost spot",
			seed:     "df55b32ff86d56703d5d0611aa5e25b13624c0c29e168b0639885650a96167d85b312bd005d560036386c480af7daccf7d9d0b99a2604485e6fc5ed493cef89b",
		},
		{
			entropy:  "df620a5757f52ae6d599d6737e820a94",
			mnemonic: "term aware noise quiz famous inflict filter depart inflict village live circle",
			seed:     "c544fa870ccee127ebb200ca2f8c2dc7fb7ba7ae3466c758adeb4850157e6aec6887a3f665dedf44aabf8618f0df600fe5f70662b888d5caba681832207545f0",
		},
		{
			entropy:  "b7080c469b1bdecd5f649ab81ef527a3",
			mnemonic: "require document balance curtain sadness grit laugh nation retreat waste enemy elevator",
			seed:     "dc49069366b49c43cc8cd5951d6e3aa30a4fd5c1a5052148558abd4e7a326d8a753f9ea97e4f308d11821e9dd958eaf09bff150c4c27756ead6522450f933563",
		},
		{
			entropy:  "9404bf555527e1225dd27cba95a11e8cc2a14fa3",
			mnemonic: "neglect chalk stem prevent lawn muffin jar exhibit ritual public element book claw pond mind",
			seed:     "69eeb23803df42a854cc57e7097f70aa8e5d1c5e79973e9af801f9e30c8cd91f0f860640aace7dcc7c401559c204ff51fed01c0eed4b0b094112fe28b6ec1ca6",
		},
		{
			entropy:  "644c63043351947638bd8d627ca08aed",
			mnemonic: "gold glide scissors grit bone deposit title random give topic cargo survey",
			seed:     "0b31a3c397edd045aa10f4dc1766665d592f74aa1144659b8091169ccd470d4432ae090f624893ee46d0b2b8233dd969dbb6d80f693ea083872e91c68a9d6d71",
		},
		{
			entropy:  "13e96f612438e0f979633b6ec2c28805",
			mnemonic: "become enter success embody mix lake tortoise guess human bid pear area",
			seed:     "ef16e3bdb61621a19f4a59cffc2563cf9bca80414df7931caccb9ae6cde3f495009560433bdaed75519555383567d8918a1323acf9f37530385221cefe66594e",
		},
		{
			entropy:  "ef710ec913d3fc199658f8848c6da7da",
			mnemonic: "urge mask rather chicken divert art floor business loyal gloom hazard remain",
			seed:     "3d1d2212d0cd07bded690041795d322108c35bfec434eb055f93dcbf732b591372222f353e71c5b1adef949007e7ed9c9a67e2bb743fb361bf33035d316ed86b",
		},
		{
			entropy:  "716ed8ae8cdc904e61eee81b763e40cb07be922c9a0bac33",
			mnemonic: "imitate item close boost simple cheese marble tackle bread rapid mother noodle know empower raven door promote odor",
			seed:     "3e329795131c02c7a3761965b9c790cc81803b879511670a2f82c8d9e9407b81546f72630721df73fd932e2b4356093cbd0a3cf931d073bea34119b9e4080fd8",
		},
		{
			entropy:  "9ff7a8c235a87bb6eb7f3b7682b000910855c4c15eee27fb",
			mnemonic: "paper run correct hero marble swarm pupil trash isolate better ability capital luxury tiny air tape child supreme",
			seed:     "ea7fffc6ff6809793cd318b2eb1561fc89a957bfc9ad833bcb6cbdc91cc0e7636d9853d84fd7c264b2122944ae3b666d7f7cfe735596aacb1d8b3ba53d0c3153",
		},
	}
}

func TestNewMnemonic(t *testing.T) {
	for i := 0; i < numTests; i++ {
		bitSize := 4 * (32 + common.RandInt()%32)
		mnemonic, err := NewMnemonic(bitSize)
		if bitSize%32 == 0 && err != nil {
			panic(err)
		}

		if bitSize%32 != 0 && err == nil {
			panic("expect an error")
		}

		if bitSize%32 != 0 && err != nil {
			continue
		}

		if err != nil {
			panic(fmt.Errorf("%v: %v", bitSize, err))
		}

		if !bip39.IsMnemonicValid(mnemonic) {
			panic(fmt.Errorf("mnemonic %v is invalid", mnemonic))
		}

		fmt.Printf("bitSize: %v, mnemonic: %v\n", bitSize, mnemonic)
	}
}

func TestNewMnemonicFromEntropy(t *testing.T) {
	for _, v := range testMnemonicVectors() {
		entropy, err := hex.DecodeString(v.entropy)
		assert.Equal(t, nil, err, fmt.Errorf("hex.DecodeString error: %v", err))

		mnemonic, err := NewMnemonicFromEntropy(entropy)
		assert.Equal(t, nil, err, fmt.Errorf("NewMnemonicFromEntropy error: %v", err))
		assert.Equal(t, v.mnemonic, mnemonic, fmt.Errorf("mnemonics for seed `%v` mismatch", v.seed))
	}
}

func TestNewMnemonicFromSeedEntropy(t *testing.T) {
	seed := common.RandBytes(32)
	entropy := common.HashB(seed)
	mnemonic, err := NewMnemonicFromEntropy(entropy[:16])
	if err != nil {
		panic(err)
	}
	fmt.Println(mnemonic)
}

func TestNewSeedFromMnemonic(t *testing.T) {
	for _, v := range testMnemonicVectors() {
		seed, err := NewSeedFromMnemonic(v.mnemonic)
		assert.Equal(t, nil, err, fmt.Errorf("NewSeedFromMnemonic error: %v", err))

		seedStr := fmt.Sprintf("%x", seed)
		assert.Equal(t, v.seed, seedStr, fmt.Errorf("seeds for mnemonic `%v` mismatch", v.mnemonic))
	}
}

func TestGetPaymentAddressV1(t *testing.T) {
	for i := 0; i < numTests; i++ {
		isNewEncoding := (common.RandInt() % 2) == 1
		privateKey := common.RandBytes(common.PrivateKeySize)
		keySet := new(key.KeySet)
		err := keySet.InitFromPrivateKeyByte(privateKey)
		assert.Equal(t, err, nil, "initKeySet returns an error: %v\n", err)

		w := new(KeyWallet)
		w.KeySet = *keySet

		PK := keySet.PaymentAddress.Pk
		TK := keySet.PaymentAddress.Tk

		paymentAddress := w.Base58CheckSerialize(PaymentAddressType)

		oldPaymentAddress, err := GetPaymentAddressV1(paymentAddress, isNewEncoding)
		assert.Equal(t, nil, err, "GetPaymentAddressV1 returns an error: %v\n", err)

		oldWallet, err := Base58CheckDeserialize(oldPaymentAddress)
		assert.Equal(t, nil, err, "deserialize returns an error: %v\n", err)

		oldPK := oldWallet.KeySet.PaymentAddress.Pk
		oldTK := oldWallet.KeySet.PaymentAddress.Tk

		assert.Equal(t, true, bytes.Equal(PK, oldPK), "public keys mismatch")
		assert.Equal(t, true, bytes.Equal(TK, oldTK), "transmission keys mismatch")
	}
}

func TestComparePaymentAddresses(t *testing.T) {
	for i := 0; i < numTests; i++ {
		privateKey := common.RandBytes(common.PrivateKeySize)
		keySet1 := new(key.KeySet)
		err := keySet1.InitFromPrivateKeyByte(privateKey)
		assert.Equal(t, err, nil, "initKeySet 1 returns an error: %v\n", err)

		keySet2 := new(key.KeySet)
		err = keySet2.InitFromPrivateKeyByte(privateKey)
		assert.Equal(t, err, nil, "initKeySet 2 returns an error: %v\n", err)

		keyWallet1 := new(KeyWallet)
		keyWallet1.KeySet = *keySet1
		keyWallet1.KeySet.PaymentAddress.OTAPublic = nil

		keyWallet2 := new(KeyWallet)
		keyWallet2.KeySet = *keySet2

		addrV1 := keyWallet1.Base58CheckSerialize(PaymentAddressType)
		assert.NotEqual(t, "", addrV1, "cannot serialize key v1")

		addrV2 := keyWallet2.Base58CheckSerialize(PaymentAddressType)
		assert.NotEqual(t, "", addrV2, "cannot serialize key v2")

		isEqual, err := ComparePaymentAddresses(addrV1, addrV2)
		assert.Equal(t, nil, err, "ComparePaymentAddresses returns an error: %v\n", err)
		assert.Equal(t, true, isEqual, "%v != %v\n", addrV1, addrV2)
	}
}

func TestGenRandomWalletForShardID(t *testing.T) {
	for i := 0; i < numTests; i++ {
		common.MaxShardNumber = common.RandInt()%7 + 1

		expectedShard := common.RandInt() % common.MaxShardNumber
		randWallet, err := GenRandomWalletForShardID(byte(expectedShard))
		assert.Equal(t, nil, err, fmt.Errorf("GenRandomWalletForShardID error: %v", err))

		pk := randWallet.KeySet.PaymentAddress.Pk
		actualShard := common.GetShardIDFromLastByte(pk[len(pk)-1])
		assert.Equal(t, expectedShard, int(actualShard), fmt.Errorf("shards mismatch with numShards = %v", common.MaxShardNumber))
	}
}
