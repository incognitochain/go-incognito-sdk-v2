package wallet

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/incognitochain/go-incognito-sdk-v2/common"

	"github.com/theedtron/go-hdwallet"

	"github.com/stretchr/testify/assert"
)

func testHDWalletVectors() []vector {
	return []vector{
		{
			mnemonic:     "search trophy awake proud sponsor toe lumber toilet sugar smoke soup joke",
			seed:         "c086086a6bb5a2e7b27d3b440e8ca849b058f8636d793e4aa366f3c59c71fa5db8793f46d87880f60290082fb40c9f43c72a93097b91aa9c6147e8214249b478",
			masterHexKey: "xprv9s21ZrQH143K2JQXJs9wvveZ7bHV2LyNcJQUiEcB5zRyuTMz8kyDCp87wHwaoXgdRoPSZ58WFpJUQ7YYSVmRobBhWpPX1poLowFREXj3UeP",
			firstKey:     "112t8rnXTLXtHoB3BY8Pnf8iHDhjGD8WoyzBAFvaKHHKU26pbZ8owQhZGhSApMqnK2KyErRZE3WYaCvm8UKjWKumoZRcpwgjAAoyyye3i2tS",
			secondKey:    "112t8rnXhFpfs7rykbQmZnp569efp5pMaEqwH1C72jUiGmD9w2Nb64i2vFqtzKg9zE2Vkc38Y4Cakbsv4ezEHU4D3zuV1EbSTVLfJMjaknCZ",
		},
		{
			mnemonic:     "any rookie avocado pipe artwork capital section draft guitar cancel swift trophy",
			seed:         "373c61d1c12e7adfed0dd2b6f9a00d80e0af6bada145536339c19eced8112a4fe2f317d597f181251ea596c3f11ae40c2f299027aea25ff6bb3a2a72b83632a4",
			masterHexKey: "xprv9s21ZrQH143K4G5Y5Pq8Z6UfWoGd3BRTTEUFjAJZYek3XjcN8yLcE4LddmNXBTeLupvc1x4KDZy2pEFwaCFAcGB2GZJ7CMmXQroz2o9hhXM",
			firstKey:     "112t8rnXTB94R1gEyPVXKa8b4R3WEa2JvxQkSfnWiUDfVikX7sxWYtQ9GXT6wRApXwnVXPtFva4W6VaeFpo4MCX9JbeWVcBgVcPh5B2EYaUP",
			secondKey:    "112t8rnXc7fLnvvZfdebZXUp1srpSLVykvVJojMdSLkpmHMKZyLvhWphysfsHYkm4juH5DvnTQJRmvM84K6N2KgVGnjcpi1fyuQEFmPvuHEJ",
		},
		{
			mnemonic:     "essay kid father raccoon frame garbage frog zoo lawn learn shove space",
			seed:         "e7ddf8c951aeea19725fa0497e3002f7d40703e12f750b6e6294fad044f491141320f3544db56ec35f5d87728bc4e98b550898e5733dfe86c5305827e9f82b20",
			masterHexKey: "xprv9s21ZrQH143K2euKznAmT8s9V286V72v8ds6jmL2FCQoX6HxaWwRW88mWUQ8cXN8DCJCTfeHooAYpqTphQNmARjGPLMvc6zYXrpRD424wDn",
			firstKey:     "112t8rnXUq9JCtYgppmgYPXNviw8HELTqVu8Sj1efCRWnDWm5gJvUC6kTzXBjoSaeka9gA1YCpWqvqyknhM8mobakBVNhGxNWuPybFJjd9P7",
			secondKey:    "112t8rnXnXuGJgCiZBiXLeiAhwhcn6CfGguGK5cqULKsxwHTPfgF8TkKxSFzWc7UJkdkrNJrEiwerWmxN13991etuNafDhabWxYjkTqhjFew",
		},
		{
			mnemonic:     "loan angry desk struggle poverty long symbol slow arrow alien point island",
			seed:         "1df25ab3435250125a4be886c8d3611b4364e110b213a978774e23b7656cda79332d5916cce3b7129ce916e71b01db562f4673b9a800ef672624d5d95b964477",
			masterHexKey: "xprv9s21ZrQH143K2CB8jNBnq5KwhCxEDM2rdoBT76TBwoVLSvqTRxfeJ8fdAmXhAU9F55bTBRWV7G5RzUyBtUUqcrVevTncL6VPtcQXqCpjaGc",
			firstKey:     "112t8rnXNY9DU9DXTzF5FHEZXnByMq1ihFWUyKH9Sk2hmSfNgsTKvKbU1o4eqyi51MLMzZPysDnYZvM48Xs3Q644gPyUpcRCkf5dr4nZoBtx",
			secondKey:    "112t8rnXhod9GHaXDj2FMHQokkBrEvkiFZsLb1gGWfSTdmmqPC6soTgKzH8R1nABSnAGWZyYGgHNfGJQmYuEh3EFdF8t9GEVweFzPhwhGw8z",
		},
		{
			mnemonic:     "wrist gym message occur hungry inquiry title flip stone piano egg person",
			seed:         "b0e90e322563cfa8ef113d0f938c07c7911eb99c2c719b81f1400903ac06a80ee8b03485f64bf01d272e2dad1fca5fee5fabf257a3d1e6f7231f1f4451fe7b6a",
			masterHexKey: "xprv9s21ZrQH143K29sG89V9Zosmx5k1c6gM2ebAB8DDriN8CD1vaXeg7WTTtQEeRmprJFGgyaJ2kecgYUEr4f8zzhZATHBUzSdQq4EunyN4HGd",
			firstKey:     "112t8rnXRWhDC4ki35NQ1xesD9kYPS9eJD7zRiZ52XT2WovjDKS9FtGnzKqQPWPMkuDNrhiDnGttf898xyEwUi7exp3boEAVyP1Z7BEiUbbP",
			secondKey:    "112t8rnXbu2TVHK1ZRia9D3odV6ao5SGvmjNg1yUPYRAghWv3uBW8VhKvuTDWpH9iTi6oEeQuvNoxbhRVXswBH2w8rb8U3fCr16QaAt9NVUP",
		},
		{
			mnemonic:     "document depth stay comfort captain marble update argue negative tent judge success",
			seed:         "9c36df1378e266ed51069646dadf6e3049eda1bc1ff4b792f4fbd4b5e81c29d11dfcfd31f45b25893a0c6271d4d324e0ebf0800554c647cb494be7c60a0a32d3",
			masterHexKey: "xprv9s21ZrQH143K4ZLzCE1ttBQxrMcM99VnTQyVMTC6MVK2NvCWTmttyYbnh8sfEWhSq2mhqewX5Y4e5ryRX6TUxk7bFbpSvzm2LoMjJctJ9t9",
			firstKey:     "112t8rnXKr8bBf2EFRncehH6wx1FD1nnLUsR84fWmVsNEVeMTJdz34B8UUEJmhXj8VAcwMkz5uxdMQKZhvTgxWZQHogTiQZ5dha3rZ2mbxRu",
			secondKey:    "112t8rnXh1b8wtanWYNwh2nbUTXcbBUC1p1CDeS2HJhma5A9wWFWjHeUw5ypAWYbRQ5hVuyGWSUee6HuJzcbQjNCzQCr5AEuRpZRUmtGpX2e",
		},
		{
			mnemonic:     "bonus deliver grass basic athlete quantum mirror climb secret depth trumpet stumble",
			seed:         "45b472b879d272e1dc21eaba15226f590d6702bbba0b52362b124c7fa39344a3d50555ad0a6b1fd774a08bd4372930def692b70a6b748f301ea18269943ffa1c",
			masterHexKey: "xprv9s21ZrQH143K2rZhHeENcEpfiDouTqQsMEojLSoRb5TEXnubRhEo1Sc414BK14EQCouW1DfECVLqXt5NNkF5SzNdDxZeSj1FytE9QWFpBRo",
			firstKey:     "112t8rnXSm9anNXkRvKt7JzefDGFLGucmraXnLSZ5gnZfc9menPAdP2g3CA2YKygw6cAGD4RzMfXBMrMB5FbgreySpAvfhbzKpaakMy4ifPW",
			secondKey:    "112t8rnXkvrtvfs5K7RggZ6uP59yw99tY8V8Q4KX9wyyAhoBMJXsjMpRRS3qz9DBseemtgs5khshJEhmaRPNLPxvgFi2dMA4yNQ9NEfsVKGV",
		},
		{
			mnemonic:     "despair body blade turn raw helmet cloud like solar lock journey close",
			seed:         "f70ed8fef90debb89a9e0fc33f148fdd5c964f6d879cea5336027ca7b3a0b3e11f9ccd965b1abc7da6bb3587539c09850a39219e31ce360b527a633282d52607",
			masterHexKey: "xprv9s21ZrQH143K3ARaNwKoJmFwub3YhwxGpTvYoaEEpgKXGgXM7uxVwR5G3qczzFWfiPSFoKc6L3WnJiN8qpjhKN5YzQag8xYXS4LNAinYBF1",
			firstKey:     "112t8rnXLdTXDVVRt83DzTAEYoYXP49gBddVztTaGGW23H3ezpYtyGiGCkph6J49u56kggkBmXLiuHgSY2pwgxvtFqYspgvPvxHbWK37xc1Q",
			secondKey:    "112t8rnXgFXFGwwbEYTs5cdUrMdZpFXAzXCaa7G57AVhLGnYoa2QAtjhMcy6N5YtRZ36qDJnaVByEf25XCe3V1CHRSByMdnZ2bvLqsfWd7VP",
		},
		{
			mnemonic:     "neglect general prize motion rebel element club misery tuition ripple kick hint",
			seed:         "e09a7ba28a49efb707e5ee6ea2871760a487a70a96b188757b72c183f20267e5a7eaeecbfb7489c181f1565e009aa8bd1cd9d79b908364d7d587e41e3efd9db6",
			masterHexKey: "xprv9s21ZrQH143K2mKf64VPYZwvabJJmKBPrU6XyWHC9MmcVgEicM6r4MFZAkokFaaPjx6Yq9d4z13dUrDASWbumbUdcyGeAcpds44iukf9efz",
			firstKey:     "112t8rnXPksrvT1LyjDTQu14gyAtxFJ7ER4HZ7jFXnGF9fzVCQY8EXNEb7kQgxmnRpd2WzAV6HKihmEQCvp7KbKcm5TdKphHsFZtHwSjQAZC",
			secondKey:    "112t8rnXbm4GQJFCTH8rgELfd1rbJtnYPkucGuuM72k9GZGhkfx7sj5imxyc7wiAMwEQw1f1LzspCamnC5QHWiSMHFaEYwomXTUjYQz8EQkc",
		},
		{
			mnemonic:     "leader dinner sister leg apology desert neck match segment scare skill foot",
			seed:         "74941057333995d3895ab4338b25a5f41575e1a2fe570b50d0fe34375c8496d21ff06abb4a16efcf70a852912a601fecb0d82ddb3a50380b68bd81c721a20263",
			masterHexKey: "xprv9s21ZrQH143K3KcgKZbhZx4JTAyGMoveV7wzTaFFUjCV5Ksmm8aoW4HWSTEUw4c8KhQFxufxSostU5VJJe4Npb6grF5TvtR2FADQMjqmwf8",
			firstKey:     "112t8rnXKQqh19XxHciDFXGVwuYvi6PKra2LsQ6Srf5T2wxqxwH4DNYzMBmsKHF1qPyCDv4AmhPBvXGdCBPWMTyq6WT94c3suf9carYjgyv6",
			secondKey:    "112t8rnXsGZNm876KvQsapg6EzuU7X6ncspcg9dqA6nmDKhHJPDP1fTy22an3dJKyCwdvLZKvLwTru6n9684w7dbjUjdk4XV6VFykJHGayqX",
		},
		{
			mnemonic:     "put fatal peanut bar bachelor mountain high mirror volume arch castle never",
			seed:         "3be8264a7b306e6cda4cb3d7f0df05d795770476edf631b80bb8b3d72752d4b3481a429fb7ebe63573b462426c9edac54d571abba90a0f790da2275535f2202f",
			masterHexKey: "xprv9s21ZrQH143K2Jo6fEumjxxDySkb5JaU7uLz3DmL394MkjjJTWVex9MR1V1C2xur76H5UynKzkojcoR42EB292nGkdF6REf2cEujmngyPna",
			firstKey:     "112t8rnXNhYJMobEL28UfUdPxQWWYWn6c4RH2bC1rKc5bxsf7f6PTTD3CV7CpwJbLYyvhPGLtkF9TLt3ZW23JobEWuw9GnXBuMNcvkQFP3mw",
			secondKey:    "112t8rnXjBdSj6wvu3f416aFUR7zTMdSPRQaGk9CuhLVr8SiFvAt7dmauZhTDs9pA2kcNLPEyPWLQkKuy22Dack3NBob34CQWA7PZBdKXQb2",
		},
		{
			mnemonic:     "protect bus fruit ensure wrestle wash kingdom vapor cage smart simple fade",
			seed:         "150ca8d7771a9255ceb4fadbcd67de853b4cae3459cb1cd955c3297065b6618dcd67b0fe4c6e4507c7f3f27bdb2da7a725a9f95a9a98cbc20028796def6368d9",
			masterHexKey: "xprv9s21ZrQH143K42xcZGwP4mgHwGz4StggqVgACo9WnmJqqekxJtiExFYiXGVxbcCyPuYRdHKsj7V2pNXngk3D4qxkuFNnfY8keKTh6UiB4LV",
			firstKey:     "112t8rnXVXee7DoAvcRW3cCR1D46AdkG4Ejukhz12wrDb9XuBhCySzPi5gqQQ23dq2KyDkHSz81K7ddXP7ENQuLux72Ds3o9C34nmt32FFrz",
			secondKey:    "112t8rnXbspcUgHMXUaSuwBDSVNTX92rchuUBsrAGWYaQrFrUjnAfrVbcC3msseu5NkuNrbkz3pDTLDj5mu7qtV2wgk6Wcexmu3YJu3TkJjK",
		},
		{
			mnemonic:     "indicate coast size wear control across frost fence door call juice inhale",
			seed:         "d6d758ff3e5c7859671baa42b20c759f1ec5fefc3a24033167ac7eb2d20ff8c3d8c1892aa495cde56df727ceff95dbe3dddda9d0802b2e41c89260fd0d85c7b3",
			masterHexKey: "xprv9s21ZrQH143K2KPF2UeW175USUz4dRXr25eSrGkZBgtF4nvB79g7W6eiDE7rpnvRmCAUjY2RDwVSdXHu13L29pGF556Z7dydfRSPTJpFpxe",
			firstKey:     "112t8rnXb3uxv12jYLSTXVrvYUD9NG33YKNzrEuF99tiJddyi3GRdQjbFh1PbznUj6uswDAAJxezaY1TWkcf38Q2Ynpo3nCPTT5cksGEovgU",
			secondKey:    "112t8rnXoFQsqUasNbPz7oQ6GBK5ZrcbuKMcTLq6PrX5WrTPFz2wJ642gwoLXTA4eTqPN7hPqnwHFq7SWzieEyHbgGFGwCiVqDmR6csdRVan",
		},
		{
			mnemonic:     "alcohol long weird stick pyramid country little legend lock priority satisfy inmate",
			seed:         "ee3e50544ac03d4896e85dbce79f24caa8d90b438d0a373a648e82e4ec7ac31df7c31fd6358c4cdc977084c2187eb9ebec5683a2f0269614ae23a22fc712c567",
			masterHexKey: "xprv9s21ZrQH143K4CtKhpGGTramHusouHuZiNzeV1WjcBgCBLNcZgkCnfsVytQ13AMJVAXVrZ5WE2qKfyKPBSU5zLQ1puZG3arxZAvgn97L6m4",
			firstKey:     "112t8rnXK3gWQFtqNCtHiS6L9VdMkriEDofDcy55uLB5G5DNfTniNYP6Rw9hHGY8HdwaEkJUB8Zf2ZqTP41aFDxaLgagCHkenvKXBad7ACDv",
			secondKey:    "112t8rnXbfh46QXeiXhMSq8eYEsgVWShiryzgcxcCMh7HLeLAhMsz9exh8XZH6jRuMZQ4EmUjKUNjt9GhkF7VPFPKADjPLwbDK3Dqwv79dPG",
		},
		{
			mnemonic:     "always uniform flight agent convince recall hard surge patient suggest eager candy",
			seed:         "b23ee1a458a9701b976d4deccaee1697d8efade98fec4730c577a6b10145aa69e46fbc36015078adc9eeb36b8938b41e7e1ca3fce936059b3f7cf1dbaee01342",
			masterHexKey: "xprv9s21ZrQH143K4QGEA2H3LfLNrkRYAfp6DozmdfJCKjgYzzF53uPJC9ZBdpzd677kGtLMbVPHi161HLKrWTg5a8X3nNE36n8W96V6vsFckdv",
			firstKey:     "112t8rnXW5ZgYBDowZQLANaitb3K22oj3ihKPRAfVWGnre7rGJitoBz4KYZ5MgSzZAcDkw4v1h51KNfx1QhqzLEKnYbnd1F7pf13J4dYVVRs",
			secondKey:    "112t8rnXbzwBoqe3XDTDpXR9vMHGk3RPtQgsusgRtUtWvsbYj66UeTQ4Pbmvz5AikTsYNXES6daCHxSvoATk1Dcy8gCYV5ydS7nADZDvcsA9",
		},
	}
}

func TestNewMasterKeyFromMnemonic(t *testing.T) {
	for _, v := range testHDWalletVectors() {
		masterWallet, err := NewMasterKeyFromMnemonic(v.mnemonic)
		if err != nil {
			panic(err)
		}

		expectedMasterWallet, err := hdwallet.StringWallet(v.masterHexKey)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, true, bytes.Equal(expectedMasterWallet.Key, masterWallet.HDKey.Key), fmt.Errorf("master keys mismatch"))
		assert.Equal(t, true, bytes.Equal(expectedMasterWallet.Chaincode, masterWallet.ChainCode), fmt.Errorf("chain codes mismatch"))

		firstWallet, err := masterWallet.DeriveChild(1)
		if err != nil {
			panic(err)
		}

		secondWallet, err := masterWallet.DeriveChild(2)
		if err != nil {
			panic(err)
		}

		firstKey := firstWallet.Base58CheckSerialize(PrivateKeyType)
		secondKey := secondWallet.Base58CheckSerialize(PrivateKeyType)

		assert.Equal(t, v.firstKey, firstKey, fmt.Sprintf("first keys mismatch: %v, %v", v.firstKey, firstKey))
		assert.Equal(t, v.secondKey, secondKey, fmt.Sprintf("second keys mismatch: %v, %v", v.secondKey, secondKey))
	}
}

func TestBase58CheckDeserialize(t *testing.T) {
	var oldWallet, newWallet *KeyWallet
	var err error
	var oldContent, newContent string

	for i := 0; i < numTests; i++ {
		common.MaxShardNumber = common.RandInt()%7 + 1
		oldWallet, err = GenRandomWalletForShardID(byte(common.RandInt() % common.MaxShardNumber))
		if err != nil {
			panic(err)
		}

		r := common.RandInt() % 4
		switch r {
		case 0:
			privateKeyStr := oldWallet.Base58CheckSerialize(PrivateKeyType)
			newWallet, err = Base58CheckDeserialize(privateKeyStr)
			if err != nil {
				panic(err)
			}

			// compare private keys
			oldContent, err = oldWallet.GetPrivateKey()
			if err != nil {
				panic(err)
			}
			newContent, err = newWallet.GetPrivateKey()
			if err != nil {
				panic(err)
			}
			if oldContent != newContent {
				panic("private keys mismatch")
			}

			// compare public keys
			oldContent, err = oldWallet.GetPublicKey()
			if err != nil {
				panic(err)
			}
			newContent, err = newWallet.GetPublicKey()
			if err != nil {
				panic(err)
			}
			if oldContent != newContent {
				panic("public keys mismatch")
			}

			// compare payment addresses
			oldContent, err = oldWallet.GetPaymentAddress()
			if err != nil {
				panic(err)
			}
			newContent, err = newWallet.GetPaymentAddress()
			if err != nil {
				panic(err)
			}
			if oldContent != newContent {
				panic("payment addresses mismatch")
			}

			// compare read-only keys
			oldContent, err = oldWallet.GetReadonlyKey()
			if err != nil {
				panic(err)
			}
			newContent, err = newWallet.GetReadonlyKey()
			if err != nil {
				panic(err)
			}
			if oldContent != newContent {
				panic("read-only keys mismatch")
			}

			// compare privateOTA keys
			oldContent, err = oldWallet.GetOTAPrivateKey()
			if err != nil {
				panic(err)
			}
			newContent, err = newWallet.GetOTAPrivateKey()
			if err != nil {
				panic(err)
			}
			if oldContent != newContent {
				panic("privateOTA keys mismatch")
			}

		case 1:
			addrStr := oldWallet.Base58CheckSerialize(PaymentAddressType)
			newWallet, err = Base58CheckDeserialize(addrStr)
			if err != nil {
				panic(err)
			}

			// compare public keys
			oldContent, err = oldWallet.GetPublicKey()
			if err != nil {
				panic(err)
			}
			newContent, err = newWallet.GetPublicKey()
			if err != nil {
				panic(err)
			}
			if oldContent != newContent {
				panic("public keys mismatch")
			}

			// compare payment addresses
			oldContent, err = oldWallet.GetPaymentAddress()
			if err != nil {
				panic(err)
			}
			newContent, err = newWallet.GetPaymentAddress()
			if err != nil {
				panic(err)
			}
			if oldContent != newContent {
				panic("payment addresses mismatch")
			}

		case 2:
			readOnlyStr := oldWallet.Base58CheckSerialize(ReadonlyKeyType)
			newWallet, err = Base58CheckDeserialize(readOnlyStr)
			if err != nil {
				panic(err)
			}

			// compare read-only keys
			oldContent, err = oldWallet.GetReadonlyKey()
			if err != nil {
				panic(err)
			}
			newContent, err = newWallet.GetReadonlyKey()
			if err != nil {
				panic(err)
			}
			if oldContent != newContent {
				panic("read-only keys mismatch")
			}

		case 3:
			otaStr := oldWallet.Base58CheckSerialize(OTAKeyType)
			newWallet, err = Base58CheckDeserialize(otaStr)
			if err != nil {
				panic(err)
			}

			// compare privateOTA keys
			oldContent, err = oldWallet.GetOTAPrivateKey()
			if err != nil {
				panic(err)
			}
			newContent, err = newWallet.GetOTAPrivateKey()
			if err != nil {
				panic(err)
			}
			if oldContent != newContent {
				panic("privateOTA keys mismatch")
			}
		}

	}
}
