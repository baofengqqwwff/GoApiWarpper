package main

import (
	"os"
	"github.com/baofengqqwwff/GoApiWarpper/binance_warpper"
	"net/http"
	"github.com/baofengqqwwff/GoApiWarpper"
	"log"
	"time"
)

func main() {
	os.Setenv("https_proxy", "http://127.0.0.1:1080")
	os.Setenv("http_proxy", "http://127.0.0.1:1080")
	binance := binance_warpper.New(http.DefaultClient,"","")
	binance.GetDepthWithWs(GoApiWarpper.NewCurrencyPair2("eth_btc"), func(depth *GoApiWarpper.Depth) {
		log.Println(depth)
	})
	time.Sleep(time.Minute)
}
