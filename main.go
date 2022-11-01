package main

import (
	"fmt"
	"io/ioutil"
)

const (
	//url = "https://zefoy.com/"
	url         = "https://froyz.com/"
	captcha_url = "a1ef290e2636bf553f39817628b6ca49.php"
	DEBUG       = true  // Basic logging
	DEBUG_2     = false // Additional logging (request responses)
)

var (
	services = map[string]string{
		"shares":    "c2VuZC9mb2xsb3dlcnNfdGlrdG9s",
		"views":     "c2VuZC9mb2xsb3dlcnNfdGlrdG9V",
		"hearts":    "c2VuZE9nb2xsb3dlcnNfdGlrdG9r",
		"followers": "c2VuZF9mb2xsb3dlcnNfdGlrdG9r",
	}
	aweme_id string
	count    uint32
)

func main() {
	banner, _ := ioutil.ReadFile("ascii.txt")
	boldRed.Printf(string(banner))
	boldRed.Printf("Enter video_id > ")
	fmt.Scanln(&aweme_id)

	go Thread("shares")
	go Thread("views")
	go Thread("hearts")
	go Thread("followers")

	select {} // Infinite "sleep" on main thread
}

func Thread(service string) {
	bot := NewBot(service)
	bot.Start()
}
