package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	Use:   "day-blocks [network] [start-block]",
	Short: "get a list of blocks that happened close to midnight local time from start block to now",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		network, ok := config.Networks[args[0]]
		if !ok {
			return fmt.Errorf("network '%s' not configured", args[0])
		}
		start, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return err
		}
		dayBlocks, err := network.GetDateBlockHeightMapping(start)
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

func (cr *NetworkDetails) GetDateBlockHeightMapping(startBlock int64) (map[time.Time]*ctypes.ResultBlock, error) {
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

	st := dates[0]
	ed := dates[len(dates)-1]

	log.Printf("finding midnight blocks for date range: start(%d/%d/%d) end(%d/%d/%d)",
		st.Month(), st.Day(), st.Year(), ed.Month(), ed.Day(), ed.Year())

	for _, date := range dates {
		if date.After(time.Now()) {
			break
		}

		estimateBlock, err := cr.GetBlock(NextBlockHeight(start, date, secondsPerBlock))
		if err != nil {
			return nil, err
		}

		secondsPerBlock = SecondsPerBlock(start, estimateBlock)

		diff := date.Sub(estimateBlock.Block.Time)

		// todo: there is an issue here where the wrong date block could get pulled. This has to do with
		// midnight height implementation below. debug and fix this.

		for math.Abs(diff.Seconds()) > 60 {
			estimateBlock, err = cr.GetBlock(NextBlockHeight(start, date, secondsPerBlock))
			if err != nil {
				return nil, err
			}
			secondsPerBlock = SecondsPerBlock(start, estimateBlock)
			diff = date.Sub(estimateBlock.Block.Time)
		}
		// TODO: do we need to set the start block = estimate block for next iteration?
		blockmap[date] = estimateBlock
	}

	log.Printf("midnight blocks identified: start(#%d) end(#%d)",
		blockmap[st].Block.Height, blockmap[ed].Block.Height)

	return blockmap, nil
}

func (cr *NetworkDetails) GetBlock(height int64) (*ctypes.ResultBlock, error) {
	node, err := cr.context.GetNode()
	if err != nil {
		return nil, err
	}
	return node.Block(context.Background(), &height)
}

func (cr *NetworkDetails) Status() (*ctypes.ResultStatus, error) {
	node, err := cr.context.GetNode()
	if err != nil {
		return nil, err
	}
	return node.Status(context.Background())
}

// round dates close to midnight (practically thats what they will be)
// up to the next date, anything small, means that we are on the right
// date and should just return those digits.

// TODO: completely review this logic
func getMidnightTime(t0 time.Time) time.Time {
	if t0.Hour() < 12 {
		return time.Date(t0.Year(), t0.Month(), t0.Day(), 0, 0, 0, 0, t0.Location())
	}
	return time.Date(t0.Year(), t0.Month(), t0.Day()-1, 0, 0, 0, 0, t0.Location())
}

func makeDates(startTime, endTime time.Time) []time.Time {
	out := []time.Time{}
	ct := startTime
	for ct.Before(endTime) {
		out = append(out, getMidnightTime(ct))
		ct = ct.Add(time.Hour * 24)
		if ct.After(time.Now()) {
			return out
		}
	}
	return out
}

func (cr *NetworkDetails) GetSecondsPerBlock(h0, h1 int64) (float64, error) {
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
