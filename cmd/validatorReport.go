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
	"fmt"
	"sort"
	"strconv"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// validatorReportCmd represents the validatorReport command
var validatorReportCmd = &cobra.Command{
	Use:   "validator-report [start-block]",
	Short: "outputs a csv of the data required for validator income reporting",
	RunE: func(cmd *cobra.Command, args []string) error {
		start, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		akashVal := NewChainReporting("akashnet-2", "http://localhost:26657", "akash", "uakt")
		akashVal.SetSDKContext()
		blocks, err := akashVal.GetDateBlockHeightMapping(start)
		if err != nil {
			return err
		}

		var eg errgroup.Group
		blockData := map[int64]AccountBlockData{}
		sem := make(chan struct{}, 50)
		blockNums := []int{}
		for k, v := range blocks {
			blockNums = append(blockNums, int(v.Block.Height))
			k, v := k, v
			eg.Go(func() error {
				bd, err := akashVal.GetBlockData(v.Block.Height, "akash1lhenngdge40r5thghzxqpsryn4x084m9jkpdg2", k)
				if err != nil {
					return err
				}
				blockData[v.Block.Height] = bd
				<-sem
				return nil
			})
			sem <- struct{}{}
		}

		sort.Ints(blockNums)
		for _, n := range blockNums {
			blockData[int64(n)].Print()
		}
		// get day blocks
		// loop over day blocks pulling data whole hog, make sure to limit concurrency
		// in a seperate routine pull the price data. be sure to build in rate limit handling
		fmt.Println(blockData)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validatorReportCmd)
}
