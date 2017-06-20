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

type BisiDriver struct{}

const siteName = "bisi"
const siteUrl = "bisi.url"
const siteCookie = "bisi.cookie"
const siteUserName = "bisi.username"
const sitePassword = "bisi.password"

func init() {
	libary.Register(siteName, BisiDriver{})
}

var (
	signed     = false
	logined    = false
	formHash   string
	qdxq       string
	iCookie    string
	curCookies []*http.Cookie = nil
)

func (d BisiDriver) Run() (rest libary.Rest, e error) {
	iLogin()
	iSign()
	rest.Code = 1
	return rest, e
}

/*初始化 formhash */
func initFormhashXq() {
	iUrl := config.GetItem(siteUrl)
	cookie := getCookie()

	urlSign := iUrl + "plugin.php?id=dsu_paulsign:sign&inajax=1"

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

	getFormHash(str)
}

/*登陆*/
func iLogin() {
	siteUrl := config.GetItem(siteUrl)
	cookie := getCookie()
	username := config.GetItem(siteUserName)
	password := config.GetItem(sitePassword)

	urlSign := siteUrl + "/member.php?mod=logging&action=login&loginsubmit=yes&infloat=yes&inajax=1"

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

	iCookie = ""
	for _, v := range curCookies {
		iCookie = iCookie + v.Name + "=" + v.Value + ";"
		//log.Println(v.Name)
		//log.Println(v.Value)
	}

	content, _ := ioutil.ReadAll(resp.Body)

	str := string(content)

	if strings.Contains(str, username) {
		log.Println("登陆成功!", "|iLogin")
		initFormhashXq()
	} else {
		log.Println("登陆失败!", "|iLogin")
	}
}

/* 签到 */
func iSign() {
	siteUrl := config.GetItem(siteUrl)
	cookie := getCookie()
	urlSign := siteUrl + "/plugin.php?id=dsu_paulsign:sign&operation=qiandao&infloat=1&inajax=1"

	param := url.Values{
		"fastreply": {"1"},
		"formhash":  {formHash},
		"qdmode":    {"1"},
		"qdxq":      {"shuai"},
		"todaysay":  {"哈哈我来签到了~"},
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


	log.Println(str)

	isSignSuccess(str)
}

func iSpeak() {
	siteUrl := config.GetItem(siteUrl)
	cookie := getCookie()
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

/* 获取 formhash */
func getFormHash(content string) {
	if strings.Contains(content, "先登錄") {
		log.Println("未登陆!", "|getFormHash")
		return
	} else {
		logined = true
	}

	reg, _ := regexp.Compile(`<input type="hidden" name="formhash" value="(.*?)">`)
	result := reg.FindStringSubmatch(content)

	if len(result) == 0 {
		log.Println("已签到!", "|getFormHash")
	} else {
		if len(result) == 2 {
			formHash = result[1]
			log.Print(formHash, "|getFormHash")
		}
	}

}

func isSigned(content string) {
	reg, _ := regexp.Compile(`<input id="(.*?)" type="radio" name="qdxq" value="(.*?)" style="display:none">`)
	result := reg.FindStringSubmatch(content)

	if len(result) == 0 {
		log.Println("已经签到!", "|isSigned")
		signed = true
	}
}

func isSignSuccess(content string) {
	if strings.Contains(content, "簽到成功") {
		log.Println("已经签到!", "|isSignSuccess")
		signed = true
	} else {
		log.Println("签到失败!", "|isSignSuccess")
	}
}

func getCookie() string {
	if len(iCookie) == 0 {
		iCookie = config.GetItem(siteCookie)
	}
	return iCookie
}
