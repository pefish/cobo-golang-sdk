package remote

import (
	"fmt"
	"github.com/pefish/cobo-golang-sdk/fixture"
	"github.com/pefish/go-logger"
	"testing"
)

func TestRemote_ListAccountBalance(t *testing.T) {
	go_logger.Logger.Init(`test`, ``)
	remote := Remote{
		BaseUrl:   `https://api.sandbox.cobo.com`,
		ApiKey:    fixture.ApiKey,
		ApiSecret: fixture.ApiSecret,
		PubKey:    fixture.Pubkey,
	}
	balances := remote.ListAccountBalance()
	for _, balance := range balances {
		fmt.Println(balance.Coin, balance.AbsBalance)
	}
}

func TestRemote_GetAccountCoinInfo(t *testing.T) {
	go_logger.Logger.Init(`test`, ``)
	remote := Remote{
		BaseUrl:   `https://api.sandbox.cobo.com`,
		ApiKey:    fixture.ApiKey,
		ApiSecret: fixture.ApiSecret,
		PubKey:    fixture.Pubkey,
	}
	balance := remote.GetAccountCoinInfo(`ETH`)
	fmt.Printf(`%#v`, balance)
}
