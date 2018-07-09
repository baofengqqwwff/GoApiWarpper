package main

import (
	"time"
	"os"
	"github.com/adshao/go-binance"
	"fmt"
)

func main() {
	os.Setenv("https_proxy", "socks5://127.0.0.1:1080")

	wsDepthHandler := func(event *binance.WsDepthEvent) {
		fmt.Println(event)
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, stopC, err := binance.WsDepthServe("LTCBTC", wsDepthHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	// use stopC to exit
	go func() {
		time.Sleep(5 * time.Second)
		stopC <- struct{}{}
	}()
	// remove this if you do not want to be blocked here
	<-doneC
}
