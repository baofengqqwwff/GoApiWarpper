package huobi_warpper

import (
	"encoding/json"
	. "github.com/baofengqqwwff/GoApiWarpper"
	"net/http"
	"github.com/nntaoli-project/GoEx/huobi"
	"github.com/nntaoli-project/GoEx"
)

var HBPOINT = NewCurrency("HBPOINT", "")

const (
	HB_POINT_ACCOUNT = "point"
	HB_SPOT_ACCOUNT  = "spot"
)

type AccountInfo struct {
	Id    string
	Type  string
	State string
}

type HuoBiPro struct {
	*huobi.HuoBiPro
}

func NewHuoBiPro(client *http.Client, apikey, secretkey, accountId string) *HuoBiPro {
	hbpro := &HuoBiPro{}
	hbpro.HuoBiPro = huobi.NewHuoBiPro(client, apikey, secretkey, accountId)

	return hbpro
}

/**
 *现货交易
 */
func NewHuoBiProSpot(client *http.Client, apikey, secretkey string) *HuoBiPro {

	hbpro := &HuoBiPro{}
	hbpro.HuoBiPro = huobi.NewHuoBiProSpot(client, apikey, secretkey)

	return hbpro
}

/**
 * 点卡账户
 */
func NewHuoBiProPoint(client *http.Client, apikey, secretkey string) *HuoBiPro {
	hbpro := &HuoBiPro{}
	hbpro.HuoBiPro = huobi.NewHuoBiProPoint(client, apikey, secretkey)

	return hbpro
}

func (hbpro *HuoBiPro) GetAccountInfo(acc string) (*AccountInfo, error) {
	goexAccountInfo, err := hbpro.HuoBiPro.GetAccountInfo(acc)
	if err != nil {
		return nil, err
	}
	goexjson, _ := json.Marshal(goexAccountInfo)
	accountInfo := &AccountInfo{}
	err = json.Unmarshal(goexjson, accountInfo)
	if err != nil {
		return nil, err
	}
	return accountInfo, nil
}

func (hbpro *HuoBiPro) GetAccount() (*Account, error) {
	goexAccountInfo, err := hbpro.HuoBiPro.GetAccount()
	if err != nil {
		return nil, err
	}
	goexjson, _ := json.Marshal(goexAccountInfo)
	accountInfo := &Account{}
	err = json.Unmarshal(goexjson, accountInfo)
	if err != nil {
		return nil, err
	}
	return accountInfo, nil
}

func (hbpro *HuoBiPro) LimitBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	goexOrder, err := hbpro.HuoBiPro.LimitBuy(amount, price, goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}
	goexjson, _ := json.Marshal(goexOrder)
	order := &Order{}
	err = json.Unmarshal(goexjson, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (hbpro *HuoBiPro) LimitSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	goexOrder, err := hbpro.HuoBiPro.LimitSell(amount, price, goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}
	goexjson, _ := json.Marshal(goexOrder)
	order := &Order{}
	err = json.Unmarshal(goexjson, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (hbpro *HuoBiPro) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	goexOrder, err := hbpro.HuoBiPro.MarketBuy(amount, price, goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}
	goexjson, _ := json.Marshal(goexOrder)
	order := &Order{}
	err = json.Unmarshal(goexjson, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (hbpro *HuoBiPro) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	goexOrder, err := hbpro.HuoBiPro.MarketSell(amount, price, goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}
	goexjson, _ := json.Marshal(goexOrder)
	order := &Order{}
	err = json.Unmarshal(goexjson, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (hbpro *HuoBiPro) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	goexOrder, err := hbpro.HuoBiPro.GetOneOrder(orderId, goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}
	goexjson, _ := json.Marshal(goexOrder)
	order := &Order{}
	err = json.Unmarshal(goexjson, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (hbpro *HuoBiPro) GetUnfinishOrders(currencyPair CurrencyPair) ([]*Order, error) {
	goexOrder, err := hbpro.HuoBiPro.GetUnfinishOrders(goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}
	goexjson, _ := json.Marshal(goexOrder)
	order := []*Order{}
	err = json.Unmarshal(goexjson, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (hbpro *HuoBiPro) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	cancelresult, err := hbpro.HuoBiPro.CancelOrder(orderId, goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	return cancelresult, err
}

func (hbpro *HuoBiPro) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]*Order, error) {
	goexOrder, err := hbpro.HuoBiPro.GetOrderHistorys(goex.NewCurrencyPair2(currency.ToSymbol("_")), currentPage, pageSize)
	if err != nil {
		return nil, err
	}
	goexjson, _ := json.Marshal(goexOrder)
	order := []*Order{}
	err = json.Unmarshal(goexjson, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (hbpro *HuoBiPro) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {
	goexTicker, err := hbpro.HuoBiPro.GetTicker(goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}
	goexjson, _ := json.Marshal(goexTicker)
	ticker := &Ticker{}
	err = json.Unmarshal(goexjson, ticker)
	if err != nil {
		return nil, err
	}
	return ticker, nil
}

func (hbpro *HuoBiPro) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {
	goexDepth, err := hbpro.HuoBiPro.GetDepth(size, goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}
	goexjson, _ := json.Marshal(goexDepth)
	depth := &Depth{}
	err = json.Unmarshal(goexjson, depth)
	if err != nil {
		return nil, err
	}
	return depth, nil
}

func (hbpro *HuoBiPro) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (hbpro *HuoBiPro) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (hbpro *HuoBiPro) GetExchangeName() (string, error) {
	return hbpro.HuoBiPro.GetExchangeName(), nil
}

func (hbpro *HuoBiPro) GetTickerWithWs(currencyPair CurrencyPair, handle func(ticker *Ticker)) error {
	goexHandler := func(goexTicker *goex.Ticker) {
		goexjson, _ := json.Marshal(goexTicker)
		warpperTicker := &Ticker{}
		err := json.Unmarshal(goexjson, warpperTicker)
		if err!=nil{
			handle(nil)
		}else{
			handle(warpperTicker)
		}

	}
	return hbpro.HuoBiPro.GetTickerWithWs(goex.NewCurrencyPair2(currencyPair.ToSymbol("_")), goexHandler)
}

func (hbpro *HuoBiPro) GetDepthWithWs(currencyPair CurrencyPair, handle func(depth *Depth)) error {
	goexHandler := func(goexDepth *goex.Depth) {
		goexjson, _ := json.Marshal(goexDepth)
		warpperDepth := &Depth{}
		err := json.Unmarshal(goexjson, warpperDepth)
		if err!=nil{
			handle(nil)
		}else{
			handle(warpperDepth)
		}

	}
	return hbpro.HuoBiPro.GetDepthWithWs(goex.NewCurrencyPair2(currencyPair.ToSymbol("_")), goexHandler)

}
