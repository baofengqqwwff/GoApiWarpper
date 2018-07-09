package GoApiWarpper

// api interface

type API interface {
	LimitBuy(amount, price string, currency CurrencyPair) (*Order, error)
	LimitSell(amount, price string, currency CurrencyPair) (*Order, error)
	MarketBuy(amount, price string, currency CurrencyPair) (*Order, error)
	MarketSell(amount, price string, currency CurrencyPair) (*Order, error)
	CancelOrder(orderId string, currency CurrencyPair) (bool, error)
	GetOneOrder(orderId string, currency CurrencyPair) (*Order, error)
	GetUnfinishOrders(currency CurrencyPair) ([]Order, error)
	GetOrderHistorys(currentPage, pageSize string,currency CurrencyPair) ([]Order, error)
	GetAccount() (*Account, error)

	GetTicker(currency CurrencyPair) (*Ticker, error)
	GetDepth(size string, currency CurrencyPair) (*Depth, error)
	GetKlineRecords(period, size, since string, currency CurrencyPair) ([]Kline, error)
	//非个人，整个交易所的交易记录
	GetTrades(since string, currencyPair CurrencyPair) ([]Trade, error)

	GetExchangeName() (string, error)
	GetDepthWithWs(currencyPair CurrencyPair, handle func(depth *Depth)) error
	CloseWs()
}
