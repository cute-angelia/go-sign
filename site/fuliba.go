package site

import (
	"go-sign/config"
	"go-sign/libary"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strconv"
	"time"
)

type FulibaDriver struct {
	signed  bool
	logined bool

	formhash string

	iCookie    string
	curCookies []*http.Cookie

	configSiteName   string
	configSiteUrl    string
	configSiteCookie string
	configSiteHost   string

	client *http.Client
}

func init() {
	cookieJar, _ := cookiejar.New(nil)

	var driver *FulibaDriver = &FulibaDriver{
		signed:           false,
		logined:          false,
		configSiteName:   "fuliba",
		configSiteUrl:    config.GetItem("fuliba.url"),
		configSiteCookie: config.GetItem("fuliba.cookie"),
		configSiteHost:   config.GetItem("fuliba.host"),
		client: &http.Client{
			Jar: cookieJar,
		},
	}
	libary.Register(driver.configSiteName, driver)
}

func (d *FulibaDriver) Run() (rest libary.Rest, e error) {

	d.iGetFormhash()
	d.iSign()

	rest.Code = 1
	return rest, e
}

func (d *FulibaDriver) getCookie() string {
	if len(d.iCookie) == 0 {
		d.iCookie = d.configSiteCookie
	}
	return d.iCookie
}

func (d *FulibaDriver) iGetFormhash() {
	t := time.Now().Unix()
	localTime := strconv.Itoa(int(t))

	url := d.configSiteUrl + "forum-2-1.html"
	cookie := d.configSiteCookie + localTime + "%09forum.php%09forumdisplay; "

	req, _ := http.NewRequest(
		"GET",
		url, nil)

	req.Header.Set("Host", d.configSiteHost)
	req.Header.Set("Origin", d.configSiteUrl)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, error := d.client.Do(req)

	if error != nil {
		log.Println(error)
	}

	content, _ := ioutil.ReadAll(resp.Body)

	str := string(content)

	pattern, _ := regexp.Compile("checkin&formhash=(.*?)&")
	result := pattern.FindStringSubmatch(str)

	if len(result) == 0 {
		log.Println("formhash not found!")
	} else {
		d.formhash = result[1]
		//log.Println(d.formhash)
	}
}

func (d *FulibaDriver) iSign() {
	t := time.Now().Unix()
	localTime := strconv.Itoa(int(t))

	url := d.configSiteUrl + "plugin.php?id=fx_checkin:checkin&formhash=" + d.formhash

	cookie := d.configSiteCookie + localTime + "%09forum.php%09forumdisplay; "

	req, _ := http.NewRequest(
		"GET",
		url, nil)

	req.Header.Set("Host", d.configSiteHost)
	req.Header.Set("Origin", d.configSiteUrl)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, error := d.client.Do(req)

	if error != nil {
		log.Println(error)
	}

	content, _ := ioutil.ReadAll(resp.Body)

	str := string(content)

	pattern, _ := regexp.Compile("<i>(.*?)</i>")
	result := pattern.FindStringSubmatch(str)

	if len(result) == 0 {
		log.Println("formhash not found!")
	} else {
		log.Println(d.configSiteName + ": 第" + result[1] + "个签到!")
	}
}
