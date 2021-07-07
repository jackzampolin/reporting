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
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// validatorReportCmd represents the validatorReport command
var validatorReportCmd = &cobra.Command{
	Use:   "validator-report [start-block]",
	Args:  cobra.ExactArgs(1),
	Short: "outputs a csv of the data required for validator income reporting",
	RunE: func(cmd *cobra.Command, args []string) error {
		start, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		akashVal := NewChainReporting("akashnet-2", "http://localhost:26657", "akash", "uakt", "akash-network")
		akashVal.SetSDKContext()
		blocks, err := akashVal.GetDateBlockHeightMapping(start)
		if err != nil {
			return err
		}

		var eg errgroup.Group
		blockData := map[int64]AccountBlockData{}
		sem := make(chan struct{}, 50)
		blockNums := []int{}
		fmt.Println("len(blocks)", len(blocks))
		for k, v := range blocks {
			blockNums = append(blockNums, int(v.Block.Height))
			k, v := k, v
			fmt.Println("starting height", v.Block.Height)
			eg.Go(func() error {
				bd, err := akashVal.GetBlockData(v.Block.Height, "akash1lhenngdge40r5thghzxqpsryn4x084m9jkpdg2", k)
				if err != nil {
					return err
				}
				price, err := akashVal.GetPrice(k)
				if err != nil {
					return err
				}
				bd.Price = price
				blockData[v.Block.Height] = bd
				<-sem
				fmt.Println("finished height", v.Block.Height)
				return nil
			})
			sem <- struct{}{}
		}

		if err := eg.Wait(); err != nil {
			return err
		}

		out := csv.NewWriter(os.Stdout)
		if err := out.Write(csvHeaders()); err != nil {
			return err
		}
		sort.Ints(blockNums)
		for _, n := range blockNums {
			if err := out.Write(blockData[int64(n)].CSVLine()); err != nil {
				return err
			}
		}
		out.Flush()
		return out.Error()
	},
}

func init() {
	rootCmd.AddCommand(validatorReportCmd)
}

func csvHeaders() []string {
	return []string{
		"date",
		"height",
		"price usd",
		"account balance native",
		"account balance usd",
		"staked balance native",
		"staked balance usd",
		"rewards balance native",
		"rewards balance usd",
		"commission balance native",
		"commission balance usd",
		"total balance native",
		"total balance usd",
	}
}
