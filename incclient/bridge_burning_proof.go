package incclient

import (
	"encoding/hex"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"math/big"
)

// BurnProof represents a proof object submitted to smart contracts for the sake of un-shielding.
type BurnProof struct {
	Instruction []byte
	Heights     [2]*big.Int

	InstPaths       [2][][32]byte
	InstPathIsLefts [2][]bool
	InstRoots       [2][32]byte
	BlkData         [2][32]byte
	SigIndices      [2][]*big.Int
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

	beaconSigIndices := make([]*big.Int, 0)
	for _, sIdx := range r.BeaconSigIndices {
		beaconSigIndices = append(beaconSigIndices, big.NewInt(int64(sIdx)))
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

	bridgeSigIndices := make([]*big.Int, 0)
	for _, sIdx := range r.BridgeSigIndices {
		bridgeSigIndices = append(bridgeSigIndices, big.NewInt(int64(sIdx)))
	}

	// Merge beacon and bridge proof
	instPaths := [2][][32]byte{beaconInstPath, bridgeInstPath}
	instPathIsLefts := [2][]bool{beaconInstPathIsLeft, bridgeInstPathIsLeft}
	instRoots := [2][32]byte{beaconInstRoot, bridgeInstRoot}
	blkData := [2][32]byte{beaconBlkData, bridgeBlkData}
	sigIndices := [2][]*big.Int{beaconSigIndices, bridgeSigIndices}
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
		SigIndices:      sigIndices,
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
		sigVs[i] = v
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
