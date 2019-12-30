package remote

import (
	"github.com/pefish/go-format"
)

type AccountBalance struct {
	Coin           string  `json:"coin"`
	DisplayCode    string  `json:"display_code"`
	Description    string  `json:"description"`
	Decimal        uint64  `json:"decimal"`
	CanDeposit     bool    `json:"can_deposit"`
	CanWithdraw    bool    `json:"can_withdraw"`
	Balance        float64 `json:"balance"`
	AbsBalance     string  `json:"abs_balance"`
	FeeCoin        string  `json:"fee_coin"`
	AbsEstimateFee string  `json:"abs_estimate_fee"`
}

func (this *Remote) ListAccountBalance() []AccountBalance {
	result := this.get(`/v1/custody/org_info/`, nil)
	assets := result.(map[string]interface{})[`assets`]
	var balances []AccountBalance
	go_format.Format.SliceToStruct(assets.([]interface{}), &balances)
	return balances
}

func (this *Remote) GetAccountCoinInfo(coin string) AccountBalance {
	result := this.get(`/v1/custody/coin_info/`, map[string]interface{}{
		`coin`: coin,
	})
	var info AccountBalance
	go_format.Format.MapToStruct(result.(map[string]interface{}), &info)
	return info
}

type FinishedTxInfo struct {
	Id                  string  `json:"id"`
	Coin                string  `json:"coin"`
	DisplayCode         string  `json:"display_code"`
	Description         string  `json:"description"`
	Decimal             uint64  `json:"decimal"`
	Address             string  `json:"address"`
	Memo                string  `json:"memo"`
	SourceAddress       string  `json:"source_address"`
	SourceAddressDetail string  `json:"source_address_detail"`
	Side                string  `json:"side"`
	Amount              string  `json:"amount"`
	AbsAmount           string  `json:"abs_amount"`
	AbsCoboFee          string  `json:"abs_cobo_fee"`
	TxId                string  `json:"txid"`
	VoutN               uint64  `json:"vout_n"`
	RequestId           *string `json:"request_id"`
	Status              string  `json:"status"`
	CreatedTime         uint64  `json:"created_time"`
	LastTime            uint64  `json:"last_time"`
	ConfirmingThreshold uint64  `json:"confirming_threshold"`
	ConfirmedNum        uint64  `json:"confirmed_num"`
	FeeCoin             string  `json:"fee_coin"`
	FeeAmount           uint64  `json:"fee_amount"`
	FeeDecimal          uint64  `json:"fee_decimal"`
	Type                string  `json:"type"`
}

func (this *Remote) GetFinishedTxInfo(id string) FinishedTxInfo {
	result := this.get(`/v1/custody/transaction/`, map[string]interface{}{
		`id`: id,
	})
	var txInfo FinishedTxInfo
	go_format.Format.MapToStruct(result.(map[string]interface{}), &txInfo)
	return txInfo
}
