package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetShardIDsFromPublicKey(t *testing.T) {
	type Tc struct {
		PubKey         []byte
		SendingShard   byte
		ReceivingShard byte
		MaxShardNumber int
	}

	testCases := []Tc{
		{[]byte{0}, 0, 0, 2},
		{[]byte{1}, 0, 1, 2},
		{[]byte{2}, 1, 0, 2},
		{[]byte{3}, 1, 1, 2},
		{[]byte{4}, 0, 0, 2},
		{[]byte{5}, 0, 1, 2},
		{[]byte{0}, 0, 0, 8},
		{[]byte{1}, 0, 1, 8},
		{[]byte{2}, 0, 2, 8},
		{[]byte{3}, 0, 3, 8},
		{[]byte{4}, 0, 4, 8},
		{[]byte{5}, 0, 5, 8},
		{[]byte{6}, 0, 6, 8},
		{[]byte{7}, 0, 7, 8},
		{[]byte{8}, 1, 0, 8},
		{[]byte{9}, 1, 1, 8},
		{[]byte{10}, 1, 2, 8},
		{[]byte{11}, 1, 3, 8},
		{[]byte{12}, 1, 4, 8},
		{[]byte{13}, 1, 5, 8},
		{[]byte{14}, 1, 6, 8},
		{[]byte{15}, 1, 7, 8},
		{[]byte{16}, 2, 0, 8},
		{[]byte{17}, 2, 1, 8},
		{[]byte{18}, 2, 2, 8},
		{[]byte{19}, 2, 3, 8},
		{[]byte{20}, 2, 4, 8},
		{[]byte{21}, 2, 5, 8},
		{[]byte{22}, 2, 6, 8},
		{[]byte{23}, 2, 7, 8},
		{[]byte{24}, 3, 0, 8},
		{[]byte{25}, 3, 1, 8},
		{[]byte{26}, 3, 2, 8},
		{[]byte{27}, 3, 3, 8},
		{[]byte{28}, 3, 4, 8},
		{[]byte{29}, 3, 5, 8},
		{[]byte{30}, 3, 6, 8},
		{[]byte{31}, 3, 7, 8},
		{[]byte{32}, 4, 0, 8},
		{[]byte{33}, 4, 1, 8},
		{[]byte{34}, 4, 2, 8},
		{[]byte{35}, 4, 3, 8},
		{[]byte{36}, 4, 4, 8},
		{[]byte{37}, 4, 5, 8},
		{[]byte{38}, 4, 6, 8},
		{[]byte{39}, 4, 7, 8},
		{[]byte{40}, 5, 0, 8},
		{[]byte{41}, 5, 1, 8},
		{[]byte{42}, 5, 2, 8},
		{[]byte{43}, 5, 3, 8},
		{[]byte{44}, 5, 4, 8},
		{[]byte{45}, 5, 5, 8},
		{[]byte{46}, 5, 6, 8},
		{[]byte{47}, 5, 7, 8},
		{[]byte{48}, 6, 0, 8},
		{[]byte{49}, 6, 1, 8},
		{[]byte{50}, 6, 2, 8},
		{[]byte{51}, 6, 3, 8},
		{[]byte{52}, 6, 4, 8},
		{[]byte{53}, 6, 5, 8},
		{[]byte{54}, 6, 6, 8},
		{[]byte{55}, 6, 7, 8},
		{[]byte{56}, 7, 0, 8},
		{[]byte{57}, 7, 1, 8},
		{[]byte{58}, 7, 2, 8},
		{[]byte{59}, 7, 3, 8},
		{[]byte{60}, 7, 4, 8},
		{[]byte{61}, 7, 5, 8},
		{[]byte{62}, 7, 6, 8},
		{[]byte{63}, 7, 7, 8},
		{[]byte{64}, 0, 0, 8},
		{[]byte{65}, 0, 1, 8},
	}

	for _, tc := range testCases {
		MaxShardNumber = tc.MaxShardNumber
		sendingShard, receivingShard := GetShardIDsFromPublicKey(tc.PubKey)
		assert.Equal(t, tc.SendingShard, sendingShard, "sendingShard failed for pubKey %v", tc.PubKey)
		assert.Equal(t, tc.ReceivingShard, receivingShard, "receivingShard failed for pubKey %v", tc.PubKey)
	}
}
