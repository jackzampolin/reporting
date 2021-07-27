package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/avast/retry-go"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	staketypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmhttp "github.com/tendermint/tendermint/rpc/client/http"
	"golang.org/x/sync/errgroup"
)

type NetworkDetails struct {
	ChainID     string `yaml:"chain-id" json:"chain-id" mapstructure:"chain-id"`
	Archive     string `yaml:"archive" json:"archive" mapstructure:"archive"`
	Prefix      string `yaml:"prefix" json:"prefix" mapstructure:"prefix"`
	Token       string `yaml:"token" json:"token" mapstructure:"token"`
	CoinGeckoID string `yaml:"coin-gecko-id" json:"coin-gecko-id" mapstructure:"coin-gecko-id"`

	context client.Context
}

type ErrRateLimitExceeded error

func (cr *NetworkDetails) GetPrice(date time.Time) (float64, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/history?date=%s&localization=false", cr.CoinGeckoID, fmt.Sprintf("%d-%d-%d", date.Day(), date.Month(), date.Year()))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/json")

	var resp *http.Response
	retry.Do(func() error {
		resp, err = http.DefaultClient.Do(req)
		switch {
		case resp.StatusCode == 429:
			return ErrRateLimitExceeded(fmt.Errorf("429"))
		case (resp.StatusCode < 200 || resp.StatusCode > 299):
			return fmt.Errorf("non 2xx or 429 status code %d", resp.StatusCode)
		case err != nil:
			return err
		default:
			return nil
		}
	}, retry.RetryIf(func(err error) bool {
		_, ok := err.(ErrRateLimitExceeded)
		return ok
	}), retry.Delay(time.Second*60))
	defer resp.Body.Close()
	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	data := priceHistory{}
	if err := json.Unmarshal(bz, &data); err != nil {
		return 0, err
	}
	return data.MarketData.CurrentPrice["usd"], nil
}

func (cr *NetworkDetails) SetSDKContext() {
	config := sdk.GetConfig()
	bech32PrefixAccAddr := cr.Prefix
	bech32PrefixAccPub := fmt.Sprintf("%spub", cr.Prefix)
	bech32PrefixValAddr := fmt.Sprintf("%svaloper", cr.Prefix)
	bech32PrefixValPub := fmt.Sprintf("%svaloperpub", cr.Prefix)
	bech32PrefixConsAddr := fmt.Sprintf("%svalcons", cr.Prefix)
	bech32PrefixConsPub := fmt.Sprintf("%svalconspub", cr.Prefix)
	config.SetBech32PrefixForAccount(bech32PrefixAccAddr, bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(bech32PrefixValAddr, bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(bech32PrefixConsAddr, bech32PrefixConsPub)
}

func (cr *NetworkDetails) Init() error {
	enc := simapp.MakeTestEncodingConfig()
	cl, err := tmhttp.New(cr.Archive, "/websocket")
	if err != nil {
		return err
	}
	cr.context = client.Context{
		Client:            cl,
		ChainID:           cr.ChainID,
		JSONMarshaler:     enc.Marshaler,
		InterfaceRegistry: enc.InterfaceRegistry,
		Input:             os.Stdin,
		Output:            os.Stdout,
		OutputFormat:      "json",
		NodeURI:           cr.Archive,
		LegacyAmino:       enc.Amino,
	}
	return nil
}

func (cr *NetworkDetails) GetBlockData(height int64, addr sdk.AccAddress, date time.Time) (AccountBlockData, error) {
	var (
		val                = sdk.ValAddress(addr)
		eg                 = errgroup.Group{}
		com, bal, rew, stk sdk.Coin
		err                error
	)
	eg.Go(func() error {
		return retry.Do(func() error {
			com, err = cr.ValidatorCommissionAtHeight(height, val)
			fmt.Println(val.String())
			fmt.Println("validator commission", err)
			return err
		})
	})
	eg.Go(func() error {
		return retry.Do(func() error {
			bal, err = cr.AccountBalanceAtHeight(height, addr)
			fmt.Println("account balance", err)
			return err
		})
	})
	eg.Go(func() error {
		return retry.Do(func() error {
			rew, err = cr.AccountRewardsAtHeight(height, addr)
			fmt.Println("rewards balance", err)
			return err
		})
	})
	eg.Go(func() error {
		return retry.Do(func() error {
			stk, err = cr.StakedTokens(height, addr)
			fmt.Println("staked tokens", err)
			return err
		})
	})
	if err := eg.Wait(); err != nil {
		return AccountBlockData{}, err
	}
	return AccountBlockData{height, bal, stk, rew, com, date, 0}, nil
}

// Raw Query Functions
func (cr *NetworkDetails) ValidatorCommissionAtHeight(height int64, val sdk.ValAddress) (sdk.Coin, error) {
	cr.context.Height = height
	res, err := distrtypes.NewQueryClient(cr.context).ValidatorCommission(context.Background(), &distrtypes.QueryValidatorCommissionRequest{ValidatorAddress: val.String()})
	if err != nil {
		return sdk.Coin{}, err
	}
	com, _ := res.Commission.Commission.TruncateDecimal()
	return sdk.NewCoin(cr.Token, com.AmountOf(cr.Token)), nil
}

func (cr *NetworkDetails) AccountRewardsAtHeight(height int64, acc sdk.AccAddress) (sdk.Coin, error) {
	cr.context.Height = height
	res, err := distrtypes.NewQueryClient(cr.context).DelegationTotalRewards(context.Background(), &distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: acc.String()})
	if err != nil {
		return sdk.Coin{}, err
	}
	rew, _ := res.Total.TruncateDecimal()
	return sdk.NewCoin(cr.Token, rew.AmountOf(cr.Token)), nil
}

func (cr *NetworkDetails) AccountBalanceAtHeight(height int64, acc sdk.AccAddress) (sdk.Coin, error) {
	cr.context.Height = height
	res, err := banktypes.NewQueryClient(cr.context).Balance(context.Background(), &banktypes.QueryBalanceRequest{Address: acc.String(), Denom: cr.Token})
	if err != nil {
		return sdk.Coin{}, err
	}
	return *res.Balance, nil
}

func (cr *NetworkDetails) StakedTokens(height int64, acc sdk.AccAddress) (sdk.Coin, error) {
	cr.context.Height = height
	res, err := staketypes.NewQueryClient(cr.context).DelegatorDelegations(context.Background(), &staketypes.QueryDelegatorDelegationsRequest{DelegatorAddr: acc.String()})
	if err != nil {
		return sdk.Coin{}, err
	}
	var tot = sdk.NewCoin(cr.Token, sdk.ZeroInt())
	for _, del := range res.DelegationResponses {
		delegation, _ := sdk.NewDecCoinFromDec(cr.Token, del.Delegation.Shares).TruncateDecimal()
		tot = tot.Add(delegation)
	}
	return tot, nil
}

type priceHistory struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Image  struct {
		Thumb string `json:"thumb"`
		Small string `json:"small"`
	} `json:"image"`
	MarketData struct {
		CurrentPrice map[string]float64 `json:"current_price"`
		MarketCap    map[string]float64 `json:"market_cap"`
		TotalVolume  map[string]float64 `json:"total_volume"`
	} `json:"market_data"`
}
