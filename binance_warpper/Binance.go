package binance_warpper

import (
	"errors"
	"fmt"
	. "github.com/baofengqqwwff/GoApiWarpper"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/copier"
	"github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/binance"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type BinanceWarpper struct {
	*binance.Binance
	ws                *WsConn
	createWsLock      sync.Mutex
	wsTickerHandleMap map[CurrencyPair]func(*Ticker)
	wsDepthHandleMap  map[CurrencyPair]func(*Depth)
}

func New(client *http.Client, api_key, secret_key string) *BinanceWarpper {
	binanceWarpper := &BinanceWarpper{}
	binanceWarpper.wsTickerHandleMap = map[CurrencyPair]func(*Ticker){}
	binanceWarpper.wsDepthHandleMap = map[CurrencyPair]func(*Depth){}
	binanceWarpper.Binance = binance.New(client, api_key, secret_key)
	return binanceWarpper
}

func (bn *BinanceWarpper) GetExchangeName() (string, error) {
	return bn.Binance.GetExchangeName(), nil
}

func (bn *BinanceWarpper) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {

	goexTicker, err := bn.Binance.GetTicker(goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}
	ticker := &Ticker{}
	err = copier.Copy(ticker, goexTicker)

	return ticker, err
}

func (bn *BinanceWarpper) GetDepth(size string, currencyPair CurrencyPair) (*Depth, error) {
	goexDepth, err := bn.Binance.GetDepth(ToInt(size), goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}

	depth := &Depth{}
	err = copier.Copy(depth, goexDepth)
	return depth, err
}

func (bn *BinanceWarpper) GetAccount() (*Account, error) {
	goexAccount, err := bn.Binance.GetAccount()
	if err != nil {
		return nil, err
	}

	account := &Account{}
	err = copier.Copy(account, goexAccount)
	return account, err
}

func (bn *BinanceWarpper) LimitBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	goexOrder, err := bn.Binance.LimitBuy(amount, price, goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}

	order := &Order{}
	err = copier.Copy(order, goexOrder)
	return order, err
}

func (bn *BinanceWarpper) LimitSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	goexOrder, err := bn.Binance.LimitSell(amount, price, goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}

	order := &Order{}
	err = copier.Copy(order, goexOrder)
	return order, err
}

func (bn *BinanceWarpper) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	goexOrder, err := bn.Binance.MarketBuy(amount, price, goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}
	order := &Order{}
	err = copier.Copy(order, goexOrder)
	return order, err
}

func (bn *BinanceWarpper) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	goexOrder, err := bn.Binance.MarketSell(amount, price, goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}
	order := &Order{}
	err = copier.Copy(order, goexOrder)
	return order, err
}

func (bn *BinanceWarpper) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	return bn.Binance.CancelOrder(orderId, goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))

}

func (bn *BinanceWarpper) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	goexOrder, err := bn.Binance.GetOneOrder(orderId, goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}
	order := &Order{}
	err = copier.Copy(order, goexOrder)
	return order, err
}

func (bn *BinanceWarpper) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	goexOrders, err := bn.Binance.GetUnfinishOrders(goex.NewCurrencyPair2(currencyPair.ToSymbol("_")))
	if err != nil {
		return nil, err
	}
	orders := []Order{}
	err = copier.Copy(orders, goexOrders)
	return orders, err
}

func (bn *BinanceWarpper) GetSymbols() ([]CurrencyPair, error) {
	exchangeInfoUri := "https://api.binance.com/api/v1/exchangeInfo"
	bodyDataMap, err := HttpGet(http.DefaultClient, exchangeInfoUri)
	if err != nil {
		return nil, err
	}
	symbols := []CurrencyPair{}
	symbolsInfo, _ := bodyDataMap["symbols"].([]interface{})
	for _, infoInterface := range symbolsInfo {
		info, _ := infoInterface.(interface{})
		symbolInfo, _ := info.(map[string]interface{})
		symbols = append(symbols, NewCurrencyPair2(symbolInfo["baseAsset"].(string)+"_"+symbolInfo["quoteAsset"].(string)))
	}
	return symbols, nil
}
func (bn *BinanceWarpper) GetKlineRecords(period, size, since string, currencyPair CurrencyPair) ([]Kline, error) {
	klineUri := "https://api.binance.com/api/v1/klines?"
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))
	if ToInt(size) > 0 {
		params.Set("limit", size)
	}
	if ToInt(since) > 0 {
		params.Set("startTime", since)
	}
	var _period string
	switch period {
	case KLINE_PERIOD_1MIN:
		{
			params.Set("interval", "1m")
			_period = "1m"
		}
	case KLINE_PERIOD_5MIN:
		{
			params.Set("interval", "5m")
			_period = "5m"
		}
	case KLINE_PERIOD_15MIN:
		{
			params.Set("interval", "15m")
			_period = "15m"
		}
	case KLINE_PERIOD_30MIN:
		{
			params.Set("interval", "30m")
			_period = "30m"
		}
	case KLINE_PERIOD_60MIN:
		{
			params.Set("interval", "1h")
			_period = "1h"
		}
	case KLINE_PERIOD_4H:
		{
			params.Set("interval", "4h")
			_period = "4h"
		}
	case KLINE_PERIOD_1DAY:
		{
			params.Set("interval", "1d")
			_period = "1d"
		}
	default:
		return nil, errors.New("do not have this period")
	}
	path := klineUri + params.Encode()
	respList, err := HttpGet3(http.DefaultClient, path, map[string]string{"X-MBX-APIKEY": ""})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	klineList := []Kline{}
	for _, resp := range respList {
		kline := Kline{}
		respBodyList, _ := resp.([]interface{})
		kline.Pair = currencyPair
		kline.Timestamp = int64(respBodyList[0].(float64))
		kline.Period = _period
		kline.Open, _ = strconv.ParseFloat(respBodyList[1].(string), 64)
		kline.High, _ = strconv.ParseFloat(respBodyList[2].(string), 64)
		kline.Low, _ = strconv.ParseFloat(respBodyList[3].(string), 64)
		kline.Close, _ = strconv.ParseFloat(respBodyList[4].(string), 64)
		kline.Vol, _ = strconv.ParseFloat(respBodyList[5].(string), 64)
		klineList = append(klineList, kline)
	}
	return klineList, nil
}

//非个人，整个交易所的交易记录
func (bn *BinanceWarpper) GetTrades(since string, currencyPair CurrencyPair) ([]Trade, error) {
	panic("not implements")
}

func (bn *BinanceWarpper) GetOrderHistorys(currentPage, pageSize string, currency CurrencyPair) ([]Order, error) {
	panic("not implements")
}

//type WsPartialDepthEvent struct {
//	Symbol       string
//	LastUpdateID int64         `json:"lastUpdateId"`
//	Bids         []DepthRecord `json:"bids"`
//	Asks         []DepthRecord `json:"asks"`
//}
//type WsHandler func(message []byte)
//type ErrHandler func(err error)

//
//func (bn *BinanceWarpper) GetDepthWithWs(currencyPair CurrencyPair, handle func(depth *Depth)) error {
//	endpoint := fmt.Sprintf("%s/%s@depth%s", "wss://stream.binance.com:9443/ws", strings.ToLower(currencyPair.ToSymbol("")), "10")
//
//	wsHandler := func(message []byte) {
//		j, err := simplejson.NewJson(message)
//		if err != nil {
//			return
//		}
//		event := new(Depth)
//		event.Pair = currencyPair
//		event.UTime = time.Now()
//		bidsLen := len(j.Get("bids").MustArray())
//		event.BidList = make([]DepthRecord, bidsLen)
//		for i := 0; i < bidsLen; i++ {
//			item := j.Get("bids").GetIndex(i)
//			price, _ := strconv.ParseFloat(item.GetIndex(0).MustString(), 64)
//			amount, _ := strconv.ParseFloat(item.GetIndex(1).MustString(), 64)
//			event.BidList[i] = DepthRecord{
//				Price:  price,
//				Amount: amount,
//			}
//		}
//		asksLen := len(j.Get("asks").MustArray())
//		event.AskList = make([]DepthRecord, asksLen)
//		for i := 0; i < asksLen; i++ {
//			item := j.Get("asks").GetIndex(i)
//			price, _ := strconv.ParseFloat(item.GetIndex(0).MustString(), 64)
//			amount, _ := strconv.ParseFloat(item.GetIndex(1).MustString(), 64)
//			event.AskList[i] = DepthRecord{
//				Price:  price,
//				Amount: amount,
//			}
//		}
//
//		handle(event)
//	}
//	return wsServe(endpoint, wsHandler)
//}
//
//func (bn *BinanceWarpper) GetDepthWithWs(currencyPair CurrencyPair, handle func(depth *Depth)) error {
//	otherPFHandler := func(event *binancews.WsPartialDepthEvent) {
//		depth := &Depth{}
//		depth.ExchangeName = "binance"
//		depth.Pair = currencyPair
//		depth.UTime = time.Now()
//		depth.AskList = make([]DepthRecord, len(event.Asks))
//		depth.BidList = make([]DepthRecord, len(event.Bids))
//		for i, ask := range event.Asks {
//			price, _ := strconv.ParseFloat(ask.Price, 64)
//			amount, _ := strconv.ParseFloat(ask.Quantity, 64)
//			depth.AskList[i] = DepthRecord{
//				Price:  price,
//				Amount: amount,
//			}
//		}
//		for i, bid := range event.Bids {
//			price, _ := strconv.ParseFloat(bid.Price, 64)
//			amount, _ := strconv.ParseFloat(bid.Quantity, 64)
//			depth.BidList[i] = DepthRecord{
//				Price:  price,
//				Amount: amount,
//			}
//		}
//
//		handle(depth)
//	}
//	errHandler := func(err error) {
//		log.Println(err)
//	}
//	_, _, err := binancews.WsPartialDepthServe(currencyPair.ToSymbol(""), "10", otherPFHandler, errHandler)
//	return err
//	//return nil
//}

func (bn *BinanceWarpper) GetDepthWithWs(currencyPair CurrencyPair, handle func(dep *Depth)) error {
	endpoint := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@depth%s", strings.ToLower(currencyPair.ToSymbol("")), "10")
	bn.createWsConn(endpoint)
	bn.ws.ReceiveMessage(func(msg []byte) {
		j, err := newJSON(msg)
		if err != nil {
			log.Println(err)
			return
		}
		event := new(WsPartialDepthEvent)
		event.Symbol = currencyPair.ToSymbol("")
		event.LastUpdateID = j.Get("lastUpdateId").MustInt64()
		bidsLen := len(j.Get("bids").MustArray())
		event.Bids = make([]Bid, bidsLen)
		for i := 0; i < bidsLen; i++ {
			item := j.Get("bids").GetIndex(i)
			event.Bids[i] = Bid{
				Price:    item.GetIndex(0).MustString(),
				Quantity: item.GetIndex(1).MustString(),
			}
		}
		asksLen := len(j.Get("asks").MustArray())
		event.Asks = make([]Ask, asksLen)
		for i := 0; i < asksLen; i++ {
			item := j.Get("asks").GetIndex(i)
			event.Asks[i] = Ask{
				Price:    item.GetIndex(0).MustString(),
				Quantity: item.GetIndex(1).MustString(),
			}
		}

		depth := &Depth{}
		depth.Exchange = "binance.com"
		depth.Pair = currencyPair
		depth.UTime = time.Now()
		depth.AskList = make([]DepthRecord, len(event.Asks))
		depth.BidList = make([]DepthRecord, len(event.Bids))
		for i, ask := range event.Asks {
			price, _ := strconv.ParseFloat(ask.Price, 64)
			amount, _ := strconv.ParseFloat(ask.Quantity, 64)
			depth.AskList[i] = DepthRecord{
				Price:  price,
				Amount: amount,
			}
		}
		for i, bid := range event.Bids {
			price, _ := strconv.ParseFloat(bid.Price, 64)
			amount, _ := strconv.ParseFloat(bid.Quantity, 64)
			depth.BidList[i] = DepthRecord{
				Price:  price,
				Amount: amount,
			}
		}
		handle(depth)

		//log.Println(string(data))
	})
	bn.wsDepthHandleMap[currencyPair] = handle
	return nil
}

func newJSON(data []byte) (j *simplejson.Json, err error) {
	j, err = simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// WsPartialDepthEvent define websocket partial depth book event
type WsPartialDepthEvent struct {
	Symbol       string
	LastUpdateID int64 `json:"lastUpdateId"`
	Bids         []Bid `json:"bids"`
	Asks         []Ask `json:"asks"`
}
type Bid struct {
	Price    string
	Quantity string
}

// Ask define ask info with price and quantity
type Ask struct {
	Price    string
	Quantity string
}

func (bn *BinanceWarpper) createWsConn(endpoint string) {
	if bn.ws == nil {
		//connect wsx
		bn.createWsLock.Lock()
		defer bn.createWsLock.Unlock()

		if bn.ws == nil {
			bn.ws = NewWsConn(endpoint)
			bn.ws.ReConnect()
			keepAlive(bn.ws, 10*time.Second)

		}
	}
}

func keepAlive(ws *WsConn, timeout time.Duration) {
	ticker := time.NewTicker(timeout)

	lastResponse := time.Now()
	ws.Conn.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		ws.UpdateActivedTime()
		return nil
	})

	go func() {
		defer ticker.Stop()
		for {
			deadline := time.Now().Add(10 * time.Second)
			err := ws.Conn.WriteControl(websocket.PingMessage, []byte{}, deadline)
			if err != nil {
				return
			}
			<-ticker.C
			if time.Now().Sub(lastResponse) > timeout {
				ws.ReConnect()

				return
			}
		}
	}()
}

func (bn *BinanceWarpper) CloseWs() {

}
