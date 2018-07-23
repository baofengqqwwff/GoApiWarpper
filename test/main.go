package main

import (
	"fmt"
	"time"
)

func main() {
	//os.Setenv("https_proxy", "http://127.0.0.1:1080")
	//os.Setenv("http_proxy", "http://127.0.0.1:1080")
	//binance := binance_warpper.New(http.DefaultClient, "", "")
	//binance.GetDepthWithWs(GoApiWarpper.NewCurrencyPair2("eth_btc"), func(depth *GoApiWarpper.Depth) {
	//	log.Println(depth)
	//})
	//time.Sleep(10*time.Minute)
	fmt.Println(time.Now().UnixNano())
	fmt.Println(time.Now().Nanosecond())
	fmt.Println(time.Now().Unix())
}
