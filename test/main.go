package main

import (
	"os"
	"net/http"
	"github.com/baofengqqwwff/GoApiWarpper"
	"log"
	"time"
	"github.com/baofengqqwwff/GoApiWarpper/huobi_warpper"
)

func main() {
	os.Setenv("https_proxy", "http://127.0.0.1:1080")
	os.Setenv("http_proxy", "http://127.0.0.1:1080")
	binance := huobi_warpper.NewHuoBiPro(http.DefaultClient, "", "","")
	binance.GetDepthWithWs(GoApiWarpper.NewCurrencyPair2("eth_btc"), func(depth *GoApiWarpper.Depth) {
		log.Println(depth)
	})
	time.Sleep(time.Minute)
}
