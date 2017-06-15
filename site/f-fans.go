package site

import (
	"go-sign/config"
	"go-sign/libary"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type FansDriver struct {
	signed     bool
	logined    bool
	formHash   string
	qdxq       string
	iCookie    string
	curCookies []*http.Cookie
	siteName     string
	siteUrl      string
	siteCookie   string
	siteUserName string
	sitePassword string
}

func init() {
	var driver *FansDriver = &FansDriver{
		signed:       false,
		logined:      false,
		siteName:     "f-fans",
		siteUrl:      "f-fans.url",
		siteCookie:   "f-fans.cookie",
		siteUserName: "f-fans.username",
		sitePassword: "f-fans.password",
	}
	libary.Register(driver.siteName, driver)
}

func (d *FansDriver) Run() (rest libary.Rest, e error) {
	d.iLogin()
	d.iSign()

	rest.Code = 1
	return rest, e
}

/*初始化 formhash */
func (d *FansDriver) initFormhashXq() {
	iUrl := config.GetItem(d.siteUrl)
	cookie := d.getCookie()

	urlSign := iUrl

	http_client := &http.Client{}
	req, _ := http.NewRequest(
		"GET",
		urlSign, nil)

	req.Header.Set("Origin", iUrl)
	req.Header.Set("Cookie", cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := http_client.Do(req)
	content, _ := ioutil.ReadAll(resp.Body)

	str := string(content)

	if strings.Contains(str, "登陆") {
		log.Println("未登陆!", "|getFormHash")
		return
	} else {
		logined = true
	}

	reg, _ := regexp.Compile(`<input type="hidden" name="formhash" value="(.*?)" />`)
	result := reg.FindStringSubmatch(str)

	if len(result) == 0 {
		log.Println("未发现formhash!", "|getFormHash")
	} else {
		if len(result) == 2 {
			d.formHash = result[1]
			log.Print(d.formHash, "|getFormHash")
		}
	}
}

/*登陆*/
func (d *FansDriver) iLogin() {
	siteUrl := config.GetItem(d.siteUrl)
	cookie := d.getCookie()
	username := config.GetItem(d.siteUserName)
	password := config.GetItem(d.sitePassword)

	urlSign := siteUrl + "member.php?mod=logging&action=login&loginsubmit=yes&infloat=yes&inajax=1"

	param := url.Values{
		"username":       {username},
		"password":       {password},
		"questionid":     {"0"},
		"answer":         {""},
		"cookietime":     {"2592000"},
		"handlekey":      {"ls"},
		"quickforward":   {"yes"},
		"fastloginfield": {"username"},
	}

	http_client := &http.Client{}
	req, _ := http.NewRequest(
		"POST",
		urlSign, strings.NewReader(param.Encode()))

	req.Header.Set("Origin", siteUrl)
	req.Header.Set("Cookie", cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := http_client.Do(req)

	curCookies = resp.Cookies()

	d.iCookie = ""
	for _, v := range curCookies {
		d.iCookie = d.iCookie + v.Name + "=" + v.Value + ";"
		//log.Println(v.Name)
		//log.Println(v.Value)
	}

	content, _ := ioutil.ReadAll(resp.Body)

	str := string(content)

	if strings.Contains(str, username) {
		log.Println("登陆成功!", "|iLogin")
		d.initFormhashXq()
	} else {
		log.Println("登陆失败!", "|iLogin")
	}
}

/* 签到 */
func (d *FansDriver) iSign() {
	siteUrl := config.GetItem(d.siteUrl)
	cookie := d.getCookie()
	urlSign := siteUrl + "plugin.php?id=k_misign:sign&operation=qiandao&inajax=1&ajaxtarget=JD_sign&formhash=" + d.formHash

	//log.Println(urlSign)

	http_client := &http.Client{}
	req, _ := http.NewRequest(
		"GET",
		urlSign, nil)

	req.Header.Set("Origin", siteUrl)
	req.Header.Set("Cookie", cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := http_client.Do(req)
	content, _ := ioutil.ReadAll(resp.Body)

	str := string(content)

	log.Println(str)

	d.isSignSuccess(str)
}

func (d *FansDriver) iSpeak() {
	siteUrl := config.GetItem(d.siteUrl)
	cookie := d.getCookie()
	urlSign := siteUrl + "/home.php?mod=spacecp&ac=doing&handlekey=doing&inajax=1"

	param := url.Values{
		"addsubmit": {"1"},
		"formhash":  {formHash},
		"referer":   {"home.php"},
		"spacenote": {"true"},
		"message":   {"哈哈我来发表心情了~"},
	}

	http_client := &http.Client{}
	req, _ := http.NewRequest(
		"POST",
		urlSign, strings.NewReader(param.Encode()))

	req.Header.Set("Origin", siteUrl)
	req.Header.Set("Cookie", cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := http_client.Do(req)
	content, _ := ioutil.ReadAll(resp.Body)

	str := string(content)

	if strings.Contains(str, "成功") {
		log.Println("发表成功!", "|iSpeak")
	} else {
		log.Println("发表失败!", "|iSpeak")
	}

}

func (d *FansDriver) isSigned(content string) {
	reg, _ := regexp.Compile(`<input id="(.*?)" type="radio" name="qdxq" value="(.*?)" style="display:none">`)
	result := reg.FindStringSubmatch(content)

	if len(result) == 0 {
		log.Println("已经签到!", "|isSigned")
		signed = true
	} else {
		//log.Print(result)
		qdxq = result[2]
	}
}

func (d *FansDriver) isSignSuccess(content string) {
	if strings.Contains(content, "簽到成功") {
		log.Println("已经签到!", "|isSignSuccess")
		signed = true
	} else {
		log.Println("签到失败!", "|isSignSuccess")
	}
}

func (d *FansDriver) getCookie() string {
	if len(d.iCookie) == 0 {
		d.iCookie = config.GetItem(d.siteCookie)
	}
	return d.iCookie
}
