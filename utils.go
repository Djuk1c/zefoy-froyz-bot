package main

import (
	b64 "encoding/base64"
	"fmt"
	"io"
	urlp "net/url"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"

	"github.com/fatih/color"
	"github.com/valyala/fasthttp"
)

var (
	boldGreen = color.New(color.FgGreen, color.Bold).SprintFunc()
	boldRed   = color.New(color.FgRed, color.Bold).SprintFunc()
	boldCyan  = color.New(color.FgCyan, color.Bold).SprintFunc()
	w         io.Writer
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
func CheckOS() {
	switch runtime.GOOS {
	case "windows":
		w = color.Output
	case "linux":
		w = os.Stdout
	}
}

func Log(msg string, clr string, service string) {
	switch clr {
	case "green":
		fmt.Fprintf(w, "[%s]: [%s]\n", color.GreenString(service), color.GreenString(msg))
	case "yellow":
		fmt.Fprintf(w, "[%s]: [%s]\n", color.YellowString(service), color.YellowString(msg))
	case "cyan":
		fmt.Fprintf(w, "[%s]: [%s]\n", boldCyan(service), boldCyan(msg))
	case "boldGreen":
		fmt.Fprintf(w, "[%s]: [%s]\n", boldGreen(service), boldGreen(msg))
	}
}

func LogErr(msg error, service string) {
	fmt.Fprintf(w, "[%s]: [%s]\n", boldRed(service), boldRed(msg))
}

func AddToCount() {
	atomic.AddUint32(&count, 1)
	fmt.Fprintf(w, "%s: %s\n", boldGreen("Successful requests"), boldGreen(strconv.Itoa(int(atomic.LoadUint32(&count)))))
}
