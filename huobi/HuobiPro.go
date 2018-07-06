package huobi

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
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
	httpClient        *http.Client
	baseUrl           string
	accountId         string
	accessKey         string
	secretKey         string
	ws                *WsConn
	createWsLock      sync.Mutex
	wsTickerHandleMap map[string]func(*Ticker)
	wsDepthHandleMap  map[string]func(*Depth)
}

func NewHuoBiPro(client *http.Client, apikey, secretkey, accountId string) *HuoBiPro {
	hbpro := new(HuoBiPro)
	hbpro.baseUrl = "https://api.huobi.br.com"
	hbpro.httpClient = client
	hbpro.accessKey = apikey
	hbpro.secretKey = secretkey
	hbpro.accountId = accountId
	hbpro.wsDepthHandleMap = make(map[string]func(*Depth))
	hbpro.wsTickerHandleMap = make(map[string]func(*Ticker))
	return hbpro
}

/**
 *现货交易
 */
func NewHuoBiProSpot(client *http.Client, apikey, secretkey string) *HuoBiPro {
	hb := NewHuoBiPro(client, apikey, secretkey, "")
	accinfo, err := hb.GetAccountInfo(HB_SPOT_ACCOUNT)
	if err != nil {
		panic(err)
	}
	hb.accountId = accinfo.Id
	log.Println("account state :", accinfo.State)
	return hb
}

/**
 * 点卡账户
 */
func NewHuoBiProPoint(client *http.Client, apikey, secretkey string) *HuoBiPro {
	hb := NewHuoBiPro(client, apikey, secretkey, "")
	accinfo, err := hb.GetAccountInfo(HB_POINT_ACCOUNT)
	if err != nil {
		panic(err)
	}
	hb.accountId = accinfo.Id
	log.Println("account state :", accinfo.State)
	return hb
}

func (hbpro *HuoBiPro) GetAccountInfo(acc string) (AccountInfo, error) {
	path := "/v1/account/accounts"
	params := &url.Values{}
	hbpro.buildPostForm("GET", path, params)

	//log.Println(hbpro.baseUrl + path + "?" + params.Encode())

	respmap, err := HttpGet(hbpro.httpClient, hbpro.baseUrl+path+"?"+params.Encode())
	if err != nil {
		return AccountInfo{}, err
	}
	//log.Println(respmap)
	if respmap["status"].(string) != "ok" {
		return AccountInfo{}, errors.New(respmap["err-code"].(string))
	}

	var info AccountInfo

	data := respmap["data"].([]interface{})
	for _, v := range data {
		iddata := v.(map[string]interface{})
		if iddata["type"].(string) == acc {
			info.Id = fmt.Sprintf("%.0f", iddata["id"])
			info.Type = acc
			info.State = iddata["state"].(string)
			break
		}
	}
	//log.Println(respmap)
	return info, nil
}

func (hbpro *HuoBiPro) GetAccount() (*Account, error) {
	path := fmt.Sprintf("/v1/account/accounts/%s/balance", hbpro.accountId)
	params := &url.Values{}
	params.Set("accountId-id", hbpro.accountId)
	hbpro.buildPostForm("GET", path, params)

	urlStr := hbpro.baseUrl + path + "?" + params.Encode()
	//println(urlStr)
	respmap, err := HttpGet(hbpro.httpClient, urlStr)

	if err != nil {
		return nil, err
	}

	//log.Println(respmap)

	if respmap["status"].(string) != "ok" {
		return nil, errors.New(respmap["err-code"].(string))
	}

	datamap := respmap["data"].(map[string]interface{})
	if datamap["state"].(string) != "working" {
		return nil, errors.New(datamap["state"].(string))
	}

	list := datamap["list"].([]interface{})
	acc := new(Account)
	acc.SubAccounts = make(map[Currency]SubAccount, 6)
	acc.Exchange = hbpro.GetExchangeName()

	subAccMap := make(map[Currency]*SubAccount)

	for _, v := range list {
		balancemap := v.(map[string]interface{})
		currencySymbol := balancemap["currency"].(string)
		currency := NewCurrency(currencySymbol, "")
		typeStr := balancemap["type"].(string)
		balance := ToFloat64(balancemap["balance"])
		if subAccMap[currency] == nil {
			subAccMap[currency] = new(SubAccount)
		}
		subAccMap[currency].Currency = currency
		switch typeStr {
		case "trade":
			subAccMap[currency].Amount = balance
		case "frozen":
			subAccMap[currency].ForzenAmount = balance
		}
	}

	for k, v := range subAccMap {
		acc.SubAccounts[k] = *v
	}

	return acc, nil
}

func (hbpro *HuoBiPro) placeOrder(amount, price string, pair CurrencyPair, orderType string) (string, error) {
	path := "/v1/order/orders/place"
	params := url.Values{}
	params.Set("account-id", hbpro.accountId)
	params.Set("amount", amount)
	params.Set("symbol", strings.ToLower(pair.ToSymbol("")))
	params.Set("type", orderType)

	switch orderType {
	case "buy-limit", "sell-limit":
		params.Set("price", price)
	}

	hbpro.buildPostForm("POST", path, &params)

	resp, err := HttpPostForm3(hbpro.httpClient, hbpro.baseUrl+path+"?"+params.Encode(), hbpro.toJson(params),
		map[string]string{"Content-Type": "application/json", "Accept-Language": "zh-cn"})
	if err != nil {
		return "", err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return "", err
	}

	if respmap["status"].(string) != "ok" {
		return "", errors.New(respmap["err-code"].(string))
	}

	return respmap["data"].(string), nil
}

func (hbpro *HuoBiPro) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := hbpro.placeOrder(amount, price, currency, "buy-limit")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		OrderID2: orderId,
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     BUY}, nil
}

func (hbpro *HuoBiPro) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := hbpro.placeOrder(amount, price, currency, "sell-limit")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		OrderID2: orderId,
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     SELL}, nil
}

func (hbpro *HuoBiPro) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := hbpro.placeOrder(amount, price, currency, "buy-market")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		OrderID2: orderId,
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     BUY_MARKET}, nil
}

func (hbpro *HuoBiPro) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := hbpro.placeOrder(amount, price, currency, "sell-market")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		OrderID2: orderId,
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     SELL_MARKET}, nil
}

func (hbpro *HuoBiPro) parseOrder(ordmap map[string]interface{}) Order {
	ord := Order{
		OrderID:    ToInt(ordmap["id"]),
		OrderID2:   fmt.Sprint(ToInt(ordmap["id"])),
		Amount:     ToFloat64(ordmap["amount"]),
		Price:      ToFloat64(ordmap["price"]),
		DealAmount: ToFloat64(ordmap["field-amount"]),
		Fee:        ToFloat64(ordmap["field-fees"]),
		OrderTime:  ToInt(ordmap["created-at"]),
	}

	state := ordmap["state"].(string)
	switch state {
	case "submitted", "pre-submitted":
		ord.Status = ORDER_UNFINISH
	case "filled":
		ord.Status = ORDER_FINISH
	case "partial-filled":
		ord.Status = ORDER_PART_FINISH
	case "canceled", "partial-canceled":
		ord.Status = ORDER_CANCEL
	default:
		ord.Status = ORDER_UNFINISH
	}

	if ord.DealAmount > 0.0 {
		ord.AvgPrice = ToFloat64(ordmap["field-cash-amount"]) / ord.DealAmount
	}

	typeS := ordmap["type"].(string)
	switch typeS {
	case "buy-limit":
		ord.Side = BUY
	case "buy-market":
		ord.Side = BUY_MARKET
	case "sell-limit":
		ord.Side = SELL
	case "sell-market":
		ord.Side = SELL_MARKET
	}
	return ord
}

func (hbpro *HuoBiPro) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	path := "/v1/order/orders/" + orderId
	params := url.Values{}
	hbpro.buildPostForm("GET", path, &params)
	respmap, err := HttpGet(hbpro.httpClient, hbpro.baseUrl+path+"?"+params.Encode())
	if err != nil {
		return nil, err
	}

	if respmap["status"].(string) != "ok" {
		return nil, errors.New(respmap["err-code"].(string))
	}

	datamap := respmap["data"].(map[string]interface{})
	order := hbpro.parseOrder(datamap)
	order.Currency = currency
	//log.Println(respmap)
	return &order, nil
}

func (hbpro *HuoBiPro) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	return hbpro.getOrders(queryOrdersParams{
		pair:   currency,
		states: "pre-submitted,submitted,partial-filled",
		size:   100,
		//direct:""
	})
}

func (hbpro *HuoBiPro) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	path := fmt.Sprintf("/v1/order/orders/%s/submitcancel", orderId)
	params := url.Values{}
	hbpro.buildPostForm("POST", path, &params)
	resp, err := HttpPostForm3(hbpro.httpClient, hbpro.baseUrl+path+"?"+params.Encode(), hbpro.toJson(params),
		map[string]string{"Content-Type": "application/json", "Accept-Language": "zh-cn"})
	if err != nil {
		return false, err
	}

	var respmap map[string]interface{}
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return false, err
	}

	if respmap["status"].(string) != "ok" {
		return false, errors.New(string(resp))
	}

	return true, nil
}

func (hbpro *HuoBiPro) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return hbpro.getOrders(queryOrdersParams{
		pair:   currency,
		size:   pageSize,
		states: "partial-canceled,filled",
		direct: "next",
	})
}

type queryOrdersParams struct {
	types,
	startDate,
	endDate,
	states,
	from,
	direct string
	size int
	pair CurrencyPair
}

func (hbpro *HuoBiPro) getOrders(queryparams queryOrdersParams) ([]Order, error) {
	path := "/v1/order/orders"
	params := url.Values{}
	params.Set("symbol", strings.ToLower(queryparams.pair.ToSymbol("")))
	params.Set("states", queryparams.states)

	if queryparams.direct != "" {
		params.Set("direct", queryparams.direct)
	}

	if queryparams.size > 0 {
		params.Set("size", fmt.Sprint(queryparams.size))
	}

	hbpro.buildPostForm("GET", path, &params)
	respmap, err := HttpGet(hbpro.httpClient, fmt.Sprintf("%s%s?%s", hbpro.baseUrl, path, params.Encode()))
	if err != nil {
		return nil, err
	}

	if respmap["status"].(string) != "ok" {
		return nil, errors.New(respmap["err-code"].(string))
	}

	datamap := respmap["data"].([]interface{})
	var orders []Order
	for _, v := range datamap {
		ordmap := v.(map[string]interface{})
		ord := hbpro.parseOrder(ordmap)
		ord.Currency = queryparams.pair
		orders = append(orders, ord)
	}

	return orders, nil
}

func (hbpro *HuoBiPro) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {
	url := hbpro.baseUrl + "/market/detail/merged?symbol=" + strings.ToLower(currencyPair.ToSymbol(""))
	respmap, err := HttpGet(hbpro.httpClient, url)
	if err != nil {
		return nil, err
	}

	if respmap["status"].(string) == "error" {
		return nil, errors.New(respmap["err-msg"].(string))
	}

	tickmap, ok := respmap["tick"].(map[string]interface{})
	if !ok {
		return nil, errors.New("tick assert error")
	}

	ticker := new(Ticker)
	ticker.Vol = ToFloat64(tickmap["amount"])
	ticker.Low = ToFloat64(tickmap["low"])
	ticker.High = ToFloat64(tickmap["high"])
	bid, isOk := tickmap["bid"].([]interface{})
	if isOk != true {
		return nil, errors.New("no bid")
	}
	ask, isOk := tickmap["ask"].([]interface{})
	if isOk != true {
		return nil, errors.New("no ask")
	}
	ticker.Buy = ToFloat64(bid[0])
	ticker.Sell = ToFloat64(ask[0])
	ticker.Last = ToFloat64(tickmap["close"])
	ticker.Date = ToUint64(respmap["ts"])

	return ticker, nil
}

func (hbpro *HuoBiPro) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	url := hbpro.baseUrl + "/market/depth?symbol=%s&type=step0"
	respmap, err := HttpGet(hbpro.httpClient, fmt.Sprintf(url, strings.ToLower(currency.ToSymbol(""))))
	if err != nil {
		return nil, err
	}

	if "ok" != respmap["status"].(string) {
		return nil, errors.New(respmap["err-msg"].(string))
	}

	tick, _ := respmap["tick"].(map[string]interface{})
	bids, _ := tick["bids"].([]interface{})
	asks, _ := tick["asks"].([]interface{})

	depth := new(Depth)
	_size := size
	for _, r := range asks {
		var dr DepthRecord
		rr := r.([]interface{})
		dr.Price = ToFloat64(rr[0])
		dr.Amount = ToFloat64(rr[1])
		depth.AskList = append(depth.AskList, dr)

		_size--
		if _size == 0 {
			break
		}
	}

	_size = size
	for _, r := range bids {
		var dr DepthRecord
		rr := r.([]interface{})
		dr.Price = ToFloat64(rr[0])
		dr.Amount = ToFloat64(rr[1])
		depth.BidList = append(depth.BidList, dr)

		_size--
		if _size == 0 {
			break
		}
	}

	sort.Sort(sort.Reverse(depth.AskList))

	return depth, nil
}

func (hbpro *HuoBiPro) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (hbpro *HuoBiPro) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (hbpro *HuoBiPro) buildPostForm(reqMethod, path string, postForm *url.Values) error {
	postForm.Set("AccessKeyId", hbpro.accessKey)
	postForm.Set("SignatureMethod", "HmacSHA256")
	postForm.Set("SignatureVersion", "2")
	postForm.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05"))
	domain := strings.Replace(hbpro.baseUrl, "https://", "", len(hbpro.baseUrl))
	payload := fmt.Sprintf("%s\n%s\n%s\n%s", reqMethod, domain, path, postForm.Encode())
	sign, _ := GetParamHmacSHA256Base64Sign(hbpro.secretKey, payload)
	postForm.Set("Signature", sign)
	return nil
}

func (hbpro *HuoBiPro) toJson(params url.Values) string {
	parammap := make(map[string]string)
	for k, v := range params {
		parammap[k] = v[0]
	}
	jsonData, _ := json.Marshal(parammap)
	return string(jsonData)
}

func (hbpro *HuoBiPro) createWsConn() {
	if hbpro.ws == nil {
		//connect wsx
		hbpro.createWsLock.Lock()
		defer hbpro.createWsLock.Unlock()

		if hbpro.ws == nil {
			hbpro.ws = NewWsConn("wss://api.huobi.br.com/ws")
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

func (hbpro *HuoBiPro) GetExchangeName() string {
	return HUOBI_PRO
}

func (hbpro *HuoBiPro) GetTickerWithWs(pair CurrencyPair, handle func(ticker *Ticker)) error {
	hbpro.createWsConn()
	sub := fmt.Sprintf("market.%s.detail", strings.ToLower(pair.ToSymbol("")))
	hbpro.wsTickerHandleMap[sub] = handle
	return hbpro.ws.Subscribe(map[string]interface{}{
		"id":  1,
		"sub": sub})
}

func (hbpro *HuoBiPro) GetDepthWithWs(pair CurrencyPair, handle func(dep *Depth)) error {
	hbpro.createWsConn()
	sub := fmt.Sprintf("market.%s.depth.step0", strings.ToLower(pair.ToSymbol("")))
	hbpro.wsDepthHandleMap[sub] = handle
	return hbpro.ws.Subscribe(map[string]interface{}{
		"id":  2,
		"sub": sub})
}
