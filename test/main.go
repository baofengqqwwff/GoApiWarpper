package main

import (
	"net/http"
	"log"
	"time"
		"sync"
	"os"
	"github.com/baofengqqwwff/GoApiWarpper/huobi_warpper"
	. "github.com/baofengqqwwff/GoApiWarpper"
)

func main() {
	os.Setenv("https_proxy" , "socks5://127.0.0.1:1080")

	huobipro := huobi_warpper.NewHuoBiPro(http.DefaultClient, "", "", "")
	huobipro.GetDepthWithWs(NewCurrencyPair2("btc_usdt"), func(depth *Depth) {
		log.Println(depth)
	})
	log.Println("启动成功")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(time.Minute)
		wg.Done()
	}()
	wg.Wait()
}
