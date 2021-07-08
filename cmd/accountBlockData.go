package cmd

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AccountBlockData struct {
	Height     int64     `json:"height"`
	Balance    sdk.Coin  `json:"balance"`
	Staked     sdk.Coin  `json:"staked"`
	Rewards    sdk.Coin  `json:"rewards"`
	Commission sdk.Coin  `json:"commission"`
	Time       time.Time `json:"time"`
	Price      float64   `json:"price"`
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

func (abd AccountBlockData) CSVLine() []string {
	return []string{
		// date
		fmt.Sprintf("%d/%d/%d", abd.Time.Month(), abd.Time.Day(), abd.Time.Year()),
		// height
		fmt.Sprintf("%d", abd.Height),
		// price usd
		fmt.Sprintf("%f", abd.Price),
		// account balance native
		abd.Balance.Amount.Quo(sdk.NewInt(1000000)).String(),
		// account balance usd
		fmt.Sprintf("%f", float64(abd.Balance.Amount.Quo(sdk.NewInt(1000000)).Int64())*abd.Price),
		// staked balance native
		abd.Staked.Amount.Quo(sdk.NewInt(1000000)).String(),
		// staked balance usd
		fmt.Sprintf("%f", float64(abd.Staked.Amount.Quo(sdk.NewInt(1000000)).Int64())*abd.Price),
		// rewards balance native
		abd.Rewards.Amount.Quo(sdk.NewInt(1000000)).String(),
		// rewards balance usd
		fmt.Sprintf("%f", float64(abd.Rewards.Amount.Quo(sdk.NewInt(1000000)).Int64())*abd.Price),
		// commission balance native
		abd.Commission.Amount.Quo(sdk.NewInt(1000000)).String(),
		// commission balance usd
		fmt.Sprintf("%f", float64(abd.Commission.Amount.Quo(sdk.NewInt(1000000)).Int64())*abd.Price),
		// total balance native
		abd.Total().Amount.Quo(sdk.NewInt(1000000)).String(),
		// total balance usd
		fmt.Sprintf("%f", float64(abd.Total().Amount.Quo(sdk.NewInt(1000000)).Int64())*abd.Price),
	}
}

func (bd AccountBlockData) Total() sdk.Coin {
	return bd.Balance.Add(bd.Staked).Add(bd.Rewards).Add(bd.Commission)
}

func (bd AccountBlockData) Print() {
	fmt.Println("balance", bd.Balance)
	fmt.Println("commission", bd.Commission)
	fmt.Println("rewards", bd.Rewards)
	fmt.Println("staked", bd.Staked)
	fmt.Println("total", bd.Total())
}
