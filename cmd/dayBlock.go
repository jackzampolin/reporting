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
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// dayBlockCmd represents the dayBlock command
var dayBlockCmd = &cobra.Command{
	Use:   "day-block [network] [height] [accaddress]",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		network, ok := config.Networks[args[0]]
		if !ok {
			return fmt.Errorf("network '%s' not configured", args[0])
		}
		network.SetSDKContext()
		address, err := sdk.AccAddressFromBech32(args[2])
		if err != nil {
			return err
		}
		start, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return err
		}

		db, err := network.GetBlockData(start, address, time.Now())
		if err != nil {
			return err
		}
		bz, err := json.Marshal(db)
		if err != nil {
			return err
		}
		fmt.Println(string(bz))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dayBlockCmd)
}
