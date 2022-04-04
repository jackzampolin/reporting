package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	rootCmd.AddCommand(validatorReportCmd)
}

// validatorReportCmd represents the validatorReport command
var validatorReportCmd = &cobra.Command{
	Use:   "validator-report [network] [validator-acc-address] [start-block]",
	Args:  cobra.ExactArgs(3),
	Short: "outputs a csv of the data required for validator income reporting",
	RunE: func(cmd *cobra.Command, args []string) error {
		network, ok := config.Networks[args[0]]
		if !ok {
			return fmt.Errorf("network '%s' not configured", args[0])
		}
		network.SetSDKContext()
		address, err := sdk.AccAddressFromBech32(args[1])
		if err != nil {
			return err
		}
		start, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			return err
		}
		blocks, err := network.GetDateBlockHeightMapping(start)
		if err != nil {
			return err
		}

		log.Printf("starting data pull, note coingecko api only returns 50 days / min...")

		var (
			eg        errgroup.Group
			blockData = map[int64]AccountBlockData{}
			sem       = make(chan struct{}, 50)
			blockNums = []int{}
			count     int
		)
		for k, v := range blocks {
			blockNums = append(blockNums, int(v.Block.Height))
			k, v := k, v
			eg.Go(func() error {
				fmt.Println("beginning pull for", k)
				bd, err := network.GetBlockData(v.Block.Height, address, k)
				if err != nil {
					fmt.Println("networkgetblock", err)
					return err
				}
				price, err := network.GetPrice(k)
				if err != nil {
					fmt.Println("getpriceerr", err)
					return err
				}
				bd.Price = price
				blockData[v.Block.Height] = bd
				<-sem
				count++
				if count%10 == 0 {
					log.Printf("%d of %d complete %f%%", count, len(blocks), (float64(count)/float64(len(blocks)))*100)
				}
				return nil
			})
			sem <- struct{}{}
		}

		// wait for all queries to return
		if err := eg.Wait(); err != nil {
			return err
		}

		// sort block numbers
		sort.Ints(blockNums)

		// create file to save csv
		file := fmt.Sprintf("report-%s-%d-%d.csv", network.ChainID, blockNums[0], blockNums[len(blockNums)-1])
		log.Printf("saving results to file: %s", file)
		out, err := os.Create(file)
		if err != nil {
			return err
		}

		// create csv writer and write data to the csv in order
		csv := csv.NewWriter(out)
		if err := csv.Write(csvHeaders()); err != nil {
			return err
		}
		for _, n := range blockNums {
			if err := csv.Write(blockData[int64(n)].CSVLine()); err != nil {
				return err
			}
		}

		// flush csv to file and return error
		csv.Flush()
		return csv.Error()
	},
}
