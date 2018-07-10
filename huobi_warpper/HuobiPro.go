package huobi_warpper

import (
	"encoding/json"
	. "github.com/baofengqqwwff/GoApiWrapper"
	"net/http"
	"github.com/nntaoli-project/GoEx/huobi"
	"github.com/nntaoli-project/GoEx"
	"sync"
	"log"
	"strings"
	"net/url"
	"errors"
	"fmt"
	"time"
	"compress/gzip"
	"bytes"
	"io/ioutil"
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
	ws                *WsConn
	createWsLock      sync.Mutex
	wsTickerHandleMap map[string]func(*Ticker)
	wsDepthHandleMap  map[string]func(*Depth)
}

func NewHuoBiPro(client *http.Client, apikey, secretkey, accountId string) *HuoBiPro {
	hbpro := &HuoBiPro{}
	hbpro.wsDepthHandleMap = make(map[string]func(*Depth))
	hbpro.wsTickerHandleMap = make(map[string]func(*Ticker))
	hbpro.HuoBiPro = huobi.NewHuoBiPro(client, apikey, secretkey, accountId)

	return hbpro
}

/**
 *现货交易
 */
func NewHuoBiProSpot(client *http.Client, apikey, secretkey string) *HuoBiPro {

	hbpro := &HuoBiPro{}
	hbpro.wsDepthHandleMap = make(map[string]func(*Depth))
	hbpro.wsTickerHandleMap = make(map[string]func(*Ticker))
	hbpro.HuoBiPro = huobi.NewHuoBiProSpot(client, apikey, secretkey)

	return hbpro
}

/**
 * 点卡账户
 */
func NewHuoBiProPoint(client *http.Client, apikey, secretkey string) *HuoBiPro {
	hbpro := &HuoBiPro{}
	hbpro.wsDepthHandleMap = make(map[string]func(*Depth))
	hbpro.wsTickerHandleMap = make(map[string]func(*Ticker))
	hbpro.HuoBiPro = huobi.NewHuoBiProPoint(client, apikey, secretkey)

	return hbpro
}

func (hbpro *HuoBiPro) GetSymbols() ([]CurrencyPair, error) {
	symbolMapInterface, err := HttpGet2(http.DefaultClient, "https://api.huobipro.com/v1/common/symbols", nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	symbolDataInterface := symbolMapInterface["data"].([]interface{})
	symbols := []CurrencyPair{}
	for _, valueInterface := range symbolDataInterface {
		value := valueInterface.(map[string]interface{})
		currencyA := value["base-currency"].(string)
		currencyB := value["quote-currency"].(string)
		symbols = append(symbols, NewCurrencyPair2(strings.ToUpper(currencyA+"_"+currencyB)))
	}
	//fmt.Println(symbolMapInterface)
	return symbols, nil
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

func (hbpro *HuoBiPro) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	goexOrder, err := hbpro.HuoBiPro.GetUnfinishOrders(goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}
	goexjson, _ := json.Marshal(goexOrder)
	order := []Order{}
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

func (hbpro *HuoBiPro) GetOrderHistorys(currentPage, pageSize string,currency CurrencyPair) ([]Order, error) {
	goexOrder, err := hbpro.HuoBiPro.GetOrderHistorys(goex.NewCurrencyPair2(currency.ToSymbol("_")), ToInt(currentPage), ToInt(pageSize))
	if err != nil {
		return nil, err
	}
	goexjson, _ := json.Marshal(goexOrder)
	order := []Order{}
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

func (hbpro *HuoBiPro) GetDepth(size string, currencyPair CurrencyPair) (*Depth, error) {
	goexDepth, err := hbpro.HuoBiPro.GetDepth(ToInt(size), goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
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

func (hbpro *HuoBiPro) GetKlineRecords(period, size, since string, currencyPair CurrencyPair) ([]Kline, error) {
	v := url.Values{}
	v.Set("symbol", strings.ToLower(currencyPair.ToSymbol("")))
	v.Set("size", size)
	var _period string
	switch period {
	case KLINE_PERIOD_1MIN:
		v.Set("period", "1min")
		_period = "1m"
	case KLINE_PERIOD_5MIN:
		v.Set("period", "5min")
		_period = "5m"
	case KLINE_PERIOD_15MIN:
		v.Set("period", "15min")
		_period = "15m"
	case KLINE_PERIOD_30MIN:
		v.Set("period", "30min")
		_period = "30m"
	case KLINE_PERIOD_60MIN:
		v.Set("period", "60min")
		_period = "60m"
	case KLINE_PERIOD_1DAY:
		v.Set("period", "1day")
		_period = "1d"
	case KLINE_PERIOD_1WEEK:
		v.Set("period", "1week")
		_period = "1w"
	default:
		return nil, errors.New("no this period")
	}
	if since != "0" {
		log.Println("please notice that since is not been used")
	}

	klineMapInterface, err := HttpGet2(http.DefaultClient, "https://api.huobipro.com/market/history/kline?"+v.Encode(), nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	klinesInterface := klineMapInterface["data"].([]interface{})
	var klineRecords []Kline

	for _, klineInterface := range klinesInterface {
		record := klineInterface.(map[string]interface{})
		klineRecord := Kline{}
		klineRecord.Close = record["close"].(float64)
		klineRecord.High = record["high"].(float64)
		klineRecord.Open = record["open"].(float64)
		klineRecord.Low = record["low"].(float64)
		klineRecord.Period = _period
		klineRecord.Pair = currencyPair
		klineRecord.Timestamp = int64(record["id"].(float64)) * 1000
		klineRecord.Vol = record["amount"].(float64)
		klineRecords = append(klineRecords, klineRecord)
	}

	var sortKlineRecords []Kline

	//按时间重新排序
	for index, _ := range klineRecords {
		sortKlineRecords = append(sortKlineRecords, klineRecords[len(klineRecords)-1-index])
	}
	return sortKlineRecords, nil
}

//非个人，整个交易所的交易记录
func (hbpro *HuoBiPro) GetTrades(since string, currencyPair CurrencyPair, ) ([]Trade, error) {
	panic("not implement")
}

func (hbpro *HuoBiPro) GetExchangeName() (string, error) {
	return hbpro.HuoBiPro.GetExchangeName(), nil
}
//
//func (hbpro *HuoBiPro) GetTickerWithWs(currencyPair CurrencyPair, handle func(ticker *Ticker)) error {
//	goexHandler := func(goexTicker *goex.Ticker) {
//		goexjson, _ := json.Marshal(goexTicker)
//		warpperTicker := &Ticker{}
//		err := json.Unmarshal(goexjson, warpperTicker)
//		if err != nil {
//			handle(nil)
//		} else {
//			handle(warpperTicker)
//		}
//
//	}
//	return hbpro.HuoBiPro.GetTickerWithWs(goex.NewCurrencyPair2(currencyPair.ToSymbol("_")), goexHandler)
//}
//
//func (hbpro *HuoBiPro) GetDepthWithWs(currencyPair CurrencyPair, handle func(depth *Depth)) error {
//	goexHandler := func(goexDepth *goex.Depth) {
//		goexjson, _ := json.Marshal(goexDepth)
//		warpperDepth := &Depth{}
//		err := json.Unmarshal(goexjson, warpperDepth)
//		if err != nil {
//			handle(nil)
//		} else {
//			handle(warpperDepth)
//		}
//
//	}
//
//	return hbpro.HuoBiPro.GetDepthWithWs(goex.NewCurrencyPair2(currencyPair.ToSymbol("_")), goexHandler)
//
//}


func (hbpro *HuoBiPro) GetTickerWithWs(pair CurrencyPair, handle func(ticker *Ticker)) error {

	hbpro.createWsConn()
	sub := fmt.Sprintf("market.%s.detail", strings.ToLower(pair.ToSymbol("")))
	hbpro.wsTickerHandleMap[sub] = handle
	return hbpro.ws.Subscribe(map[string]interface{}{
		"id":  1,
		"sub": sub})
}
func (hbpro *HuoBiPro) CloseWs(){
	hbpro.ws.CloseWs()
}

func (hbpro *HuoBiPro) GetDepthWithWs(pair CurrencyPair, handle func(dep *Depth)) error {
	hbpro.createWsConn()
	sub := fmt.Sprintf("market.%s.depth.step0", strings.ToLower(pair.ToSymbol("")))
	hbpro.wsDepthHandleMap[sub] = handle
	return hbpro.ws.Subscribe(map[string]interface{}{
		"id":  2,
		"sub": sub})
}


func (hbpro *HuoBiPro) createWsConn() {
	if hbpro.ws == nil {
		//connect wsx
		hbpro.createWsLock.Lock()
		defer hbpro.createWsLock.Unlock()

		if hbpro.ws == nil {
			hbpro.ws = NewWsConn("wss://api.huobipro.com/ws")
			hbpro.ws.Heartbeat(func() interface{} {
				return map[string]interface{}{
					"ping": time.Now().Unix()}
			}, 5*time.Second)
			hbpro.ws.ReConnect()
			hbpro.ws.ReceiveMessage(func(msg []byte) {
				gzipreader, _ := gzip.NewReader(bytes.NewReader(msg))
				data, _ := ioutil.ReadAll(gzipreader)
				datamap := make(map[string]interface{})
				err := json.Unmarshal(data, &datamap)
				if err != nil {
					log.Println("json unmarshal error for ", string(data))
					return
				}

				if datamap["ping"] != nil {
					//log.Println(datamap)
					hbpro.ws.UpdateActivedTime()
					hbpro.ws.WriteJSON(map[string]interface{}{
						"pong": datamap["ping"]}) // 回应心跳
					return
				}

				if datamap["pong"] != nil { //
					hbpro.ws.UpdateActivedTime()
					return
				}

				if datamap["id"] != nil { //忽略订阅成功的回执消息
					log.Println(string(data))
					return
				}

				ch, isok := datamap["ch"].(string)
				if !isok {
					log.Println("error:", string(data))
					return
				}

				tick := datamap["tick"].(map[string]interface{})
				if hbpro.wsTickerHandleMap[ch] != nil {
					return
				}

				if hbpro.wsDepthHandleMap[ch] != nil {
					(hbpro.wsDepthHandleMap[ch])(hbpro.parseDepthData(tick))
					return
				}

				//log.Println(string(data))
			})
		}
	}
}

func (hbpro *HuoBiPro) parseDepthData(tick map[string]interface{}) *Depth {
	bids, _ := tick["bids"].([]interface{})
	asks, _ := tick["asks"].([]interface{})

	depth := new(Depth)
	for _, r := range asks {
		var dr DepthRecord
		rr := r.([]interface{})
		dr.Price = ToFloat64(rr[0])
		dr.Amount = ToFloat64(rr[1])
		depth.AskList = append(depth.AskList, dr)
	}

	for _, r := range bids {
		var dr DepthRecord
		rr := r.([]interface{})
		dr.Price = ToFloat64(rr[0])
		dr.Amount = ToFloat64(rr[1])
		depth.BidList = append(depth.BidList, dr)
	}

	return depth
}
