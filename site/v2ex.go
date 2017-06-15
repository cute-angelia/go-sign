package site

import (
	"go-sign/config"
	"go-sign/libary"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

type V2exDriver struct {
	signed       bool
	logined      bool
	iCookie      string
	curCookies   []*http.Cookie
	siteName     string
	siteUrl      string
	siteCookie   string
	siteUserName string
	sitePassword string

	userNameKey string
	passKey     string
	once        string

	client *http.Client
}

func init() {
	cookieJar, _ := cookiejar.New(nil)

	var driver *V2exDriver = &V2exDriver{
		signed:       false,
		logined:      false,
		siteName:     "v2ex",
		siteUrl:      "v2ex.url",
		siteCookie:   "v2ex.cookie",
		siteUserName: "v2ex.username",
		sitePassword: "v2ex.password",
		client: &http.Client{
			Jar: cookieJar,
		},
	}
	libary.Register(driver.siteName, driver)
}

func (d *V2exDriver) Run() (rest libary.Rest, e error) {

	d.iGetLogin()
	d.iLogin()
	d.iSign()

	rest.Code = 1
	return rest, e
}

func (d *V2exDriver) getCookie() string {
	if len(d.iCookie) == 0 {
		d.iCookie = config.GetItem(d.siteCookie)
	}
	return d.iCookie
}

func (d *V2exDriver) iGetLogin() {
	siteUrl := config.GetItem(d.siteUrl)

	urlSign := siteUrl + "signin"

	req, _ := http.NewRequest(
		"GET",
		urlSign, nil)

	req.Header.Set("Host", "www.v2ex.com")
	req.Header.Set("Origin", "https://www.v2ex.com")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, error := d.client.Do(req)

	if error != nil {
		log.Println(error)
	}

	content, _ := ioutil.ReadAll(resp.Body)

	str := string(content)

	d.getUsernameKey(str)
	d.getPassKey(str)
	d.getOnce(str)
}

func (d *V2exDriver) iLogin() {
	siteUrl := config.GetItem(d.siteUrl)
	username := config.GetItem(d.siteUserName)
	password := config.GetItem(d.sitePassword)

	urlSign := siteUrl + "signin"

	param := url.Values{
		d.userNameKey: {username},
		d.passKey:     {password},
		"once":        {d.once},
		"next":        {"/"},
	}

	req, _ := http.NewRequest(
		"POST",
		urlSign, strings.NewReader(param.Encode()))

	req.Header.Set("Host", "www.v2ex.com")
	req.Header.Set("Origin", "https://www.v2ex.com")
	req.Header.Set("Referer", "https://www.v2ex.com/signin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")

	resp, error := d.client.Do(req)

	if error != nil {
		log.Println(error)
	}

	curCookies = resp.Cookies()

	d.iCookie = ""
	for _, v := range curCookies {
		d.iCookie = d.iCookie + v.Name + "=" + v.Value + ";"
	}

	content, _ := ioutil.ReadAll(resp.Body)

	str := string(content)

	if strings.Contains(str, "signout") {
		log.Println("登陆成功!", "|iLogin")
	} else {
		log.Println("登陆失败!", "|iLogin")
	}
}

func (d *V2exDriver) iSign() {
	siteUrl := config.GetItem(d.siteUrl)
	cookie := d.getCookie()
	urlSign := siteUrl + "mission/daily"

	req, _ := http.NewRequest(
		"GET",
		urlSign, nil)

	req.Header.Set("Host", "www.v2ex.com")
	req.Header.Set("Origin", "https://www.v2ex.com")
	req.Header.Set("Cookie", cookie)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")

	resp, error := d.client.Do(req)

	if error != nil {
		log.Println(error)
	}

	content, _ := ioutil.ReadAll(resp.Body)

	str := string(content)

	d.getDailyOnce(str)

	// xxxx
	url2 := siteUrl + "mission/daily/redeem?once=" + d.once
	req2, _ := http.NewRequest(
		"GET",
		url2, nil)
	req.Header.Set("Host", "www.v2ex.com")
	req.Header.Set("Origin", "https://www.v2ex.com")
	req.Header.Set("Cookie", cookie)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
	resp2, _ := d.client.Do(req2)
	content2, _ := ioutil.ReadAll(resp2.Body)
	str2 := string(content2)

	if strings.Contains(str2, "已成功领取每日登录奖励") {
		log.Println("领取成功!")
	} else {
		log.Println("领取失败!")
	}
}

func (d *V2exDriver) getUsernameKey(content string) {
	reg, _ := regexp.Compile(`<input type="text" class="sl" name="(.*?)"`)
	result := reg.FindStringSubmatch(content)

	if len(result) == 0 {
		log.Println("not found!")
	} else {
		d.userNameKey = result[1]
	}
}

func (d *V2exDriver) getPassKey(content string) {
	reg, _ := regexp.Compile(`<input type="password" class="sl" name="(.*?)"`)
	result := reg.FindStringSubmatch(content)

	if len(result) == 0 {
		log.Println("not found!")
	} else {
		d.passKey = result[1]
	}
}

func (d *V2exDriver) getOnce(content string) {
	reg, _ := regexp.Compile(`<input type="hidden" value="(.*?)" name="once" />`)
	result := reg.FindStringSubmatch(content)

	if len(result) == 0 {
		log.Println("not found!")
	} else {
		d.once = result[1]
	}
}

func (d *V2exDriver) getDailyOnce(content string) {
	reg, _ := regexp.Compile(`<input type="button" class="super normal button" value="领取 X 铜币" onclick="location.href = '/mission/daily/redeem\?once=(.*?)';" />`)
	result := reg.FindStringSubmatch(content)

	if len(result) == 0 {
		log.Println("not found!")
	} else {
		d.once = result[1]
	}
}
