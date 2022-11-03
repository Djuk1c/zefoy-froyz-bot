package main

import (
	b64 "encoding/base64"
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/anaskhan96/soup"
	"github.com/valyala/fasthttp"
)

const (
	CONN_TIMEOUT = 30 * time.Second
)

type Bot struct {
	client  *fasthttp.Client
	sessid  string
	captcha string
	alpha   string
	beta    string
	service string
}

func NewBot(service string) *Bot {
	bot := new(Bot)
	bot.client = &fasthttp.Client{}
	bot.service = service
	return bot
}

func (b *Bot) GetSessionID() (validProxy bool) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.SetRequestURI(url)
	SetHeaders(req)
	err := b.client.DoDeadline(req, resp, time.Now().Add(CONN_TIMEOUT))
	if err != nil {
		if DEBUG {
			LogErr(err, b.service)
		}
		return false
	}
	sessid := string(resp.Header.PeekCookie("PHPSESSID"))
	if sessid == "" {
		if DEBUG {
			LogErr(errors.New("PHPSESSID not found."), b.service)
		}
		return false
	}
	sessid = sessid[:strings.IndexByte(sessid, ';')]
	b.sessid = sessid
	if DEBUG {
		Log(sessid, "green", b.service)
	}

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	return true
}

func (b *Bot) GetCaptcha() {
	// Get image
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.SetRequestURI(url)
	SetHeaders(req, b.sessid)
	unix := strconv.FormatInt(time.Now().Unix(), 10)
	num := strconv.FormatUint(rand.Uint64(), 10)[:8]
	captchaUrl := url + captcha_url + "?_CAPTCHA&t=0." + num + "+" + unix
	req.SetRequestURI(captchaUrl)
	b.client.DoDeadline(req, resp, time.Now().Add(CONN_TIMEOUT))

	// Some guys OCR API (TODO: Implement own ocr w gosserect or something)
	imgCode := b64.StdEncoding.EncodeToString(resp.Body())
	json := gabs.New()
	json.Set(string(imgCode), "img")
	creq := fasthttp.AcquireRequest()
	cresp := fasthttp.AcquireResponse()
	creq.SetRequestURI("https://api.sandroputraa.com/zefoy.php")
	creq.Header.SetMethod("POST")
	creq.Header.Set("Content-Type", "application/json")
	creq.Header.Set("Auth", "sandrocods")
	creq.Header.Set("Host", "api.sandroputraa.com")
	creq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")
	creq.SetBodyString(json.String())
	fasthttp.Do(creq, cresp)
	resJson, _ := gabs.ParseJSON(cresp.Body())
	captcha := resJson.Search("Data").String()
	captcha = strings.ToLower(captcha[1 : len(captcha)-1])
	b.captcha = captcha
	if DEBUG {
		Log(captcha, "green", b.service)
	}

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	defer fasthttp.ReleaseRequest(creq)
	defer fasthttp.ReleaseResponse(cresp)
}

func (b *Bot) GetAlphaKey() (valid bool) {
	// Can panic if captcha result is wrong
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/x-www-form-urlencoded")
	req.SetBodyString("captcha_secure=" + b.captcha + "&r75619cf53f5a5d7ba6af82edfec3bf0=")
	SetHeaders(req, b.sessid)
	b.client.DoDeadline(req, resp, time.Now().Add(CONN_TIMEOUT))

	if DEBUG_2 {
		Log("ALPHA_KEY_RESPONSE: "+string(resp.Body()), "yellow", b.service)
	}
	if !strings.Contains(string(resp.Body()), "name") { // Bad gateway/empty response
		LogErr(errors.New("Zefoy.com empty/bad response, waiting 15 seconds"), b.service)
		return false
	}
	doc := soup.HTMLParse(string(resp.Body()))
	alpha := doc.Find("input").Attrs()["name"] // Panics here when captcha wrong (handled)
	b.alpha = alpha
	if DEBUG {
		Log("ALPHA_KEY: "+alpha, "green", b.service)
	}

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	return true
}

func (b *Bot) GetBetaKey() (timer int) {
	// Will panic if on timer
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.SetRequestURI(url + services[b.service])
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/x-www-form-urlencoded")
	req.SetBodyString(b.alpha + "=https://www.tiktok.com/@djukicdev/video/" + aweme_id)
	SetHeaders(req, b.sessid)
	b.client.DoDeadline(req, resp, time.Now().Add(CONN_TIMEOUT))

	decode := Decode(string(resp.Body()))
	doc := soup.HTMLParse(decode)
	if DEBUG_2 {
		Log("BETA_KEY_RESPONSE:\n"+decode, "yellow", b.service)
	}
	if strings.Contains(decode, "This service is currently not working") {
		LogErr(errors.New("This service is currently disabled."), b.service)
		return -2
	} else if strings.Contains(decode, "Too many requests. Please slow down") {
		LogErr(errors.New("Too many requests, waiting 15 seconds."), b.service)
		return 15
	} else if strings.Contains(decode, "Session expired. Please re-login") {
		LogErr(errors.New("Session expired, waiting 15 seconds"), b.service)
		return -1
	} else if strings.Contains(decode, "Server too busy.") || decode == "" {
		LogErr(errors.New("Server busy, waiting 15 seconds"), b.service)
		return 15
	} else if strings.Contains(decode, "Checking Timer") { // Means we're on timer
		temp := strings.Split(decode, "\n")[3]
		re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		temp = re.FindAllString(temp, -1)[0]
		timer, _ := strconv.Atoi(temp)
		Log("BETA_TIMER: "+strconv.Itoa(timer), "boldGreen", b.service)
		return timer
	}
	beta := doc.Find("input").Attrs()["name"]
	b.beta = beta
	if DEBUG {
		Log("BETA_KEY: "+beta, "green", b.service)
	}

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	return 0
}

func (b *Bot) Submit() (timer int) {
	// Can respond with Timed out, Success, or Time Left
	// View loop can later "panic" if phpsessid has expired
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURI(url + services[b.service])
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/x-www-form-urlencoded")
	req.SetBodyString(b.beta + "=" + aweme_id)
	SetHeaders(req, b.sessid)

	b.client.DoDeadline(req, resp, time.Now().Add(CONN_TIMEOUT))
	decode := Decode(string(resp.Body()))
	if DEBUG_2 {
		Log("SUBMIT_RESPONSE: "+decode, "yellow", b.service)
	}
	if strings.Contains(decode, "Too many requests. Please slow down") {
		LogErr(errors.New("Too many requests, waiting 15 seconds."), b.service)
		return 15
	} else if decode == "" {
		LogErr(errors.New("Empty BETA Response, waiting 15 seconds"), b.service)
		return 15
	} else if strings.Contains(decode, "sent") {
		Log("SENT", "cyan", b.service)
		AddToCount()
	} else if strings.Contains(decode, "ltm=") { // Means we're on timer
		// Sometimes ltm=939; ends up here
		temp := strings.Split(decode, "\n")[2] // 2?
		re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		temp = re.FindAllString(temp, -1)[0]
		timer, _ := strconv.Atoi(temp)
		Log("BETA_TIMER: "+strconv.Itoa(timer), "boldGreen", b.service)
		return timer
	}

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	return 0
}

func (b *Bot) Start() {
	// Hacky but it works ¯\_(ツ)_/¯
	timer := 0
	if !b.GetSessionID() { //PHPSESS Not found
		return
	}
	b.GetCaptcha()
	if !b.GetAlphaKey() { // Zefoy responded empty, IP banned/Captcha wrong
		time.Sleep(15 * time.Second)
		b.Start()
	}

	for {
		time.Sleep(time.Duration(timer) * time.Second)
		timer = 0

		t := b.GetBetaKey()
		if t == -2 { // Service disabled
			return
		}
		if t == -1 { // Session expired
			time.Sleep(15 * time.Second)
			b.Start()
		}
		if t != 0 {
			timer = t
			continue
		}
		t = b.Submit()
		if t != 0 {
			timer = t
			continue
		}
	}
}
