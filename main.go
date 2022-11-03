package main

import (
	"fmt"
	"io/ioutil"
	neturl "net/url"
	"os"
	"strconv"
	"strings"

	screen "github.com/aditya43/clear-shell-screen-golang"
)

const (
	url = "https://zefoy.com/"
	//url         = "https://froyz.com/"
	captcha_url = "a1ef290e2636bf553f39817628b6ca49.php"
)

var (
	DEBUG    = true  // Basic logging
	DEBUG_2  = false // Additional logging (request responses)
	services = map[string]string{
		"views":     "c2VuZC9mb2xsb3dlcnNfdGlrdG9V",
		"hearts":    "c2VuZE9nb2xsb3dlcnNfdGlrdG9r",
		"followers": "c2VuZF9mb2xsb3dlcnNfdGlrdG9r",
		"favorites": "c2VuZF9mb2xsb3dlcnNfdGlrdG9L",
		"shares":    "c2VuZC9mb2xsb3dlcnNfdGlrdG9s",
	}
	aweme_id string
	count    uint32
)

func main() {
	screen.Clear()
	CheckArguments()
	banner, _ := ioutil.ReadFile("ascii.txt")
	boldRed.Printf(string(banner))
	boldRed.Printf("Enter URL/VideoID > ")
	fmt.Scanln(&aweme_id)
	aweme_id = ProcessUrl(aweme_id)

	for i, _ := range services {
		go Thread(i)
	}

	select {} // Infinite "sleep" on main thread
}

func Thread(service string) {
	bot := NewBot(service)
	bot.Start()
}

func CheckArguments() {
	if len(os.Args) > 1 {
		if os.Args[1] == "--debug" {
			DEBUG_2 = true
		}
	}
}

func ProcessUrl(link string) string {
	if _, err := strconv.Atoi(link); err == nil {
		return link
	} else {
		u, _ := neturl.Parse(link)
		return strings.Split(u.Path, "/")[3]
	}
}
