package main

import (
	b64 "encoding/base64"
	"fmt"
	urlp "net/url"
	"sync/atomic"

	"github.com/fatih/color"
	"github.com/valyala/fasthttp"
)

var (
	boldGreen = color.New(color.FgGreen, color.Bold)
	boldRed   = color.New(color.FgRed, color.Bold)
	boldCyan  = color.New(color.FgCyan, color.Bold)
)

func ReverseString(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}
func Decode(str string) (result string) {
	rev, _ := urlp.QueryUnescape(ReverseString(string(str)))
	decode, _ := b64.StdEncoding.DecodeString(rev)
	return string(decode)
}

func SetHeaders(req *fasthttp.Request, cookie ...string) {
	req.Header.Set("origin", url)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)           Chrome/101.0.4951.54 Safari/537.36")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("Host", "zefoy.com")
	if len(cookie) > 0 {
		req.Header.Set("cookie", cookie[0])
	}
}

func Log(msg string, clr string, service string) {
	switch clr {
	case "green":
		color.Green("[%s]: [%s]", service, msg)
	case "yellow":
		color.Yellow("[%s]: [%s]", service, msg)
	case "cyan":
		boldCyan.Printf("[%s]: [%s]\n", service, msg)
	case "boldGreen":
		boldGreen.Printf("[%s]: [%s]\n", service, msg)
	}
}

func LogErr(msg error, service string) {
	fmt.Printf("\033[31m[%s] : [%s]\n\033[0m", service, msg)
}

func AddToCount() {
	atomic.AddUint32(&count, 1)
	fmt.Printf("\033]0;Successful requests: %d\007", atomic.LoadUint32(&count))
}
