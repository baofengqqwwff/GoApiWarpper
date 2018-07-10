package GoApiWarpper

// api interface

type API interface {
	LimitBuy(amount, price string, currencyPair CurrencyPair) (*Order, error)
	LimitSell(amount, price string, currencyPair CurrencyPair) (*Order, error)
	MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error)
	MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error)
	CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error)
	GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error)
	GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error)
	GetOrderHistorys(currentPage, pageSize string,currency CurrencyPair) ([]Order, error)
	GetAccount() (*Account, error)

	GetTicker(currencyPair CurrencyPair) (*Ticker, error)
	GetDepth(size string, currencyPair CurrencyPair) (*Depth, error)
	GetKlineRecords(period, size, since string, currencyPair CurrencyPair) ([]Kline, error)
	//非个人，整个交易所的交易记录
	GetTrades(since string, currencyPair CurrencyPair) ([]Trade, error)

	GetExchangeName() (string, error)
	GetDepthWithWs(currencyPair CurrencyPair, handle func(depth *Depth)) error
	CloseWs()
}
