package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/dgryski/go-farm"
	"github.com/lithammer/go-jump-consistent-hash"
)

type FarmHash struct {
	buf bytes.Buffer
}

func (f *FarmHash) Write(p []byte) (n int, err error) {
	return f.buf.Write(p)
}

func (f *FarmHash) Reset() {
	f.buf.Reset()
}

func (f *FarmHash) Sum64() uint64 {
	// https://github.com/dgryski/go-farm
	return farm.Hash64(f.buf.Bytes())
}

func randRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func numberPercentage(i int, total int) int {
	return i * 100 / total
}

func maxRepetitions(m map[int]int64) int64 {
	max := int64(0)

	for _, val := range m {
		if val > max {
			max = val
		}
	}

	return max
}

func main() {
	rand.Seed(time.Now().UnixNano())

	itemsCount := 5000000
	offset := time.Second * 57
	nodesCount := int32(5)

	shardsDistribution := make(map[int32]int)
	keysToShard := make([]string, 0, itemsCount)
	keys := make(map[int]int64)

	now := time.Now()
	for i := 0; i < itemsCount; i++ {
		key := randRange(int(now.UnixMilli()), int(now.Add(offset).UnixMilli()))
		keys[key]++

		keysToShard = append(keysToShard, fmt.Sprintf("%d", key))
	}

	for _, item := range keysToShard {
		node := jump.HashString(item, nodesCount, jump.NewCRC64())
		shardsDistribution[node]++
	}

	fmt.Println("done")

	fmt.Println(shardsDistribution)
	// fmt.Println(keys)

	for nodeID, cnt := range shardsDistribution {
		fmt.Printf("Node: %d has got %d%%\n", nodeID, numberPercentage(cnt, itemsCount))
	}

	fmt.Println("Unique keys: ", len(keys))
	fmt.Println("Max items with the same key: ", maxRepetitions(keys))
}
