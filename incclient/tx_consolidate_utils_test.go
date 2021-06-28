package incclient

import (
	"fmt"
	"testing"
)

func TestEstimateNumTxs(t *testing.T) {
	testData := [][]int{
		{1, 1, 0},
		{10, 1, 1},
		{31, 1, 2},
		{30, 1, 1},
		{900, 1, 31},
		{901, 1, 32},
	}

	for i, data := range testData {
		numTxs := estimateNumTxs(data[0], data[1])
		if numTxs != data[2] {
			panic(fmt.Sprintf("i = %v, data %v, numTxs %v\n", i, data, numTxs))
		}
	}
}
