package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	staketypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmhttp "github.com/tendermint/tendermint/rpc/client/http"
	"golang.org/x/sync/errgroup"
)

type ChainReporting struct {
	ChainID string
	NodeURI string
	Prefix  string
	Token   string
	Context client.Context
}

func (cr *ChainReporting) SetSDKContext() {
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

func NewChainReporting(chainid, nodeuri, prefix, token string) *ChainReporting {
	enc := simapp.MakeTestEncodingConfig()
	cr := &ChainReporting{
		ChainID: chainid,
		NodeURI: nodeuri,
		Prefix:  prefix,
		Token:   token,
	}
	cl, err := tmhttp.New(cr.NodeURI, "/websocket")
	if err != nil {
		panic(err)
	}
	cr.Context = client.Context{
		Client:            cl,
		ChainID:           cr.ChainID,
		JSONMarshaler:     enc.Marshaler,
		InterfaceRegistry: enc.InterfaceRegistry,
		Input:             os.Stdin,
		Output:            os.Stdout,
		OutputFormat:      "json",
		NodeURI:           cr.NodeURI,
		LegacyAmino:       enc.Amino,
	}
	return cr
}

// func main() {
// 	GetAkashBlock(168228)
// }

func GetAkashBlock(height int64, blockTime time.Time) {
	akashVal := NewChainReporting("akashnet-2", "http://localhost:26657", "akash", "uakt")
	akashVal.SetSDKContext()
	bd, err := akashVal.GetBlockData(height, "akash1lhenngdge40r5thghzxqpsryn4x084m9jkpdg2", blockTime)
	if err != nil {
		log.Fatal(err)
	}
	bd.Print()
}

type AccountBlockData struct {
	Height     int64
	Balance    sdk.Coin
	Staked     sdk.Coin
	Rewards    sdk.Coin
	Commission sdk.Coin
	Time       time.Time
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

func (cr *ChainReporting) GetBlockData(height int64, valoper string, date time.Time) (AccountBlockData, error) {
	addr, err := sdk.AccAddressFromBech32(valoper)
	if err != nil {
		return AccountBlockData{}, err
	}
	val := sdk.ValAddress(addr)
	var com, bal, rew, stk sdk.Coin
	var eg = errgroup.Group{}
	eg.Go(func() error {
		com, err = cr.ValidatorCommissionAtHeight(height, val)
		return err
	})
	eg.Go(func() error {
		bal, err = cr.AccountBalanceAtHeight(height, addr)
		return err
	})
	eg.Go(func() error {
		rew, err = cr.AccountRewardsAtHeight(height, addr)
		return err
	})
	eg.Go(func() error {
		stk, err = cr.StakedTokens(height, addr)
		return err
	})
	if err := eg.Wait(); err != nil {
		return AccountBlockData{}, err
	}
	return AccountBlockData{height, bal, stk, rew, com, date}, nil
}

func (cr *ChainReporting) ValidatorCommissionAtHeight(height int64, val sdk.ValAddress) (sdk.Coin, error) {
	cr.Context.Height = height
	res, err := distrtypes.NewQueryClient(cr.Context).ValidatorCommission(context.Background(), &distrtypes.QueryValidatorCommissionRequest{ValidatorAddress: val.String()})
	if err != nil {
		return sdk.Coin{}, err
	}
	com, _ := res.Commission.Commission.TruncateDecimal()
	return sdk.NewCoin(cr.Token, com.AmountOf(cr.Token)), nil
}

func (cr *ChainReporting) AccountRewardsAtHeight(height int64, acc sdk.AccAddress) (sdk.Coin, error) {
	cr.Context.Height = height
	res, err := distrtypes.NewQueryClient(cr.Context).DelegationTotalRewards(context.Background(), &distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: acc.String()})
	if err != nil {
		return sdk.Coin{}, err
	}
	rew, _ := res.Total.TruncateDecimal()
	return sdk.NewCoin(cr.Token, rew.AmountOf(cr.Token)), nil
}

func (cr *ChainReporting) AccountBalanceAtHeight(height int64, acc sdk.AccAddress) (sdk.Coin, error) {
	cr.Context.Height = height
	res, err := banktypes.NewQueryClient(cr.Context).Balance(context.Background(), &banktypes.QueryBalanceRequest{Address: acc.String(), Denom: cr.Token})
	if err != nil {
		return sdk.Coin{}, err
	}
	return *res.Balance, nil
}

func (cr *ChainReporting) StakedTokens(height int64, acc sdk.AccAddress) (sdk.Coin, error) {
	cr.Context.Height = height
	res, err := staketypes.NewQueryClient(cr.Context).DelegatorDelegations(context.Background(), &staketypes.QueryDelegatorDelegationsRequest{DelegatorAddr: acc.String()})
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
