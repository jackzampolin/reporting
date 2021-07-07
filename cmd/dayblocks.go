/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func init() {
	rootCmd.AddCommand(dayblocksCmd)
}

// dayblocksCmd represents the dayblocks command
var dayblocksCmd = &cobra.Command{
	Use:   "day-blocks [start-block]",
	Short: "get a list of blocks that happened close to midnight local time from start block to now",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		akashVal := NewChainReporting("akashnet-2", "http://localhost:26657", "akash", "uakt")
		dayBlocks, err := akashVal.GetDateBlockHeightMapping(start)
		if err != nil {
			return err
		}
		bz, _ := json.Marshal(dayBlocksHuman(dayBlocks))
		fmt.Println(string(sdk.MustSortJSON(bz)))
		return nil
	},
}

func dayBlocksHuman(db map[time.Time]*ctypes.ResultBlock) map[string]int64 {
	out := map[string]int64{}
	for k, v := range db {
		out[k.String()] = v.Block.Height
	}
	return out
}

func (cr *ChainReporting) GetDateBlockHeightMapping(startBlock int64) (map[time.Time]*ctypes.ResultBlock, error) {
	start, err := cr.GetBlock(startBlock)
	if err != nil {
		return nil, err
	}

	status, err := cr.Status()
	if err != nil {
		return nil, err
	}

	var (
		blockmap        = map[time.Time]*ctypes.ResultBlock{}
		secondsPerBlock = start.Block.Time.Sub(status.SyncInfo.LatestBlockTime).Seconds() / float64(startBlock-status.SyncInfo.LatestBlockHeight)
		dates           = makeDates(start.Block.Time, status.SyncInfo.LatestBlockTime)
	)

	for _, date := range dates {
		estimateBlock, err := cr.GetBlock(NextBlockHeight(start, date, secondsPerBlock))
		if err != nil {
			return nil, err
		}
		secondsPerBlock = SecondsPerBlock(start, estimateBlock)

		diff := date.Sub(estimateBlock.Block.Time)
		for math.Abs(diff.Seconds()) > 60 {
			estimateBlock, err = cr.GetBlock(NextBlockHeight(start, date, secondsPerBlock))
			if err != nil {
				return nil, err
			}
			secondsPerBlock = SecondsPerBlock(start, estimateBlock)
			diff = date.Sub(estimateBlock.Block.Time)
		}
		blockmap[date] = estimateBlock
	}

	return blockmap, nil
}

func (cr *ChainReporting) GetBlock(height int64) (*ctypes.ResultBlock, error) {
	node, err := cr.Context.GetNode()
	if err != nil {
		return nil, err
	}
	return node.Block(context.Background(), &height)
}

func (cr *ChainReporting) Status() (*ctypes.ResultStatus, error) {
	node, err := cr.Context.GetNode()
	if err != nil {
		return nil, err
	}
	return node.Status(context.Background())
}

func getMidnightTime(t0 time.Time) time.Time {
	return time.Date(t0.Year(), t0.Month(), t0.Day()+1, 0, 0, 0, 0, t0.Location())
}

func makeDates(startTime, endTime time.Time) []time.Time {
	out := []time.Time{}
	ct := startTime
	for ct.Before(endTime) {
		out = append(out, getMidnightTime(ct))
		ct = ct.Add(time.Hour * 24)
	}
	return out
}

func (cr *ChainReporting) GetSecondsPerBlock(h0, h1 int64) (float64, error) {
	b0, err := cr.GetBlock(h0)
	if err != nil {
		return 0, err
	}
	b1, err := cr.GetBlock(h1)
	if err != nil {
		return 0, err
	}
	return SecondsPerBlock(b0, b1), nil
}

func SecondsPerBlock(b0, b1 *ctypes.ResultBlock) float64 {
	return b0.Block.Time.Sub(b1.Block.Time).Seconds() / float64(b0.Block.Height-b1.Block.Height)
}

func NextBlockHeight(startBlock *ctypes.ResultBlock, nextDate time.Time, secondsPerBlock float64) int64 {
	return startBlock.Block.Height + int64(nextDate.Sub(startBlock.Block.Time).Seconds()/secondsPerBlock)
}
