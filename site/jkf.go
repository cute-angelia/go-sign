package site

import (
	"go-sign/config"
	"go-sign/libary"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"regexp"
)

type JkfDriver struct {
	Cookie   string
	Title    string
	Url      string
	FormHash string
}

func init() {
	cookie := config.GetItem("jkf.cookie")
	title := config.GetItem("jkf.title")
	iurl := config.GetItem("jkf.url")

	libary.Register("jkf", JkfDriver{
		Cookie: cookie,
		Title:  title,
		Url:    iurl,
	})
}

func (d JkfDriver) Run() (rest libary.Rest, e error) {
	d.initFormhashXq();
	d.iSign()
	rest.Code = 1
	return rest, e
}

/*初始化 fromhash */
func (d *JkfDriver) initFormhashXq() {

	urlSign := d.Url + "/plugin.php?id=dsu_paulsign:sign&inajax=1"

	http_client := &http.Client{}
	req, _ := http.NewRequest(
		"GET",
		urlSign, nil)

	req.Header.Set("Origin", d.Url)
	req.Header.Set("Cookie", d.Cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := http_client.Do(req)
	content, _ := ioutil.ReadAll(resp.Body)

	str := string(content)
	d.getFormHash(str)
}

/* 获取 formhash */
func (d *JkfDriver) getFormHash(content string) string {
	if strings.Contains(content, "先登錄") {
		log.Println("未登陆!", "|getFormHash")
		return ""
	} else {
		logined = true
	}

	reg, _ := regexp.Compile(`<input type="hidden" name="formhash" value="(.*?)">`)
	result := reg.FindStringSubmatch(content)

	if len(result) == 0 {
		log.Println("已签到!", "|getFormHash")
	} else {
		if len(result) == 2 {
			d.FormHash = result[1]
			log.Print(d.FormHash, "|getFormHash")
		}
	}
	return d.FormHash
}

/* 签到 */
func (d JkfDriver) iSign() {
	urlSign := d.Url + "/plugin.php?id=dsu_paulsign:sign&operation=qiandao&infloat=1&inajax=1"
	// https://www.jkforum.net/plugin/?id=straightdisplay:ajax&ajaxact=straightdisplay&formhash=591f4c96&inajax=1&ajaxtarget=forumlist

	param := url.Values{
		"fastreply": {"1"},
		"formhash":  {d.FormHash},
		"qdmode":    {"1"},
		"qdxq":      {"shuai"},
		"todaysay":  {"哈哈我来签到了~"},
	}

	http_client := &http.Client{}
	req, _ := http.NewRequest(
		"POST",
		urlSign, strings.NewReader(param.Encode()))

	req.Header.Set("Origin", d.Url)
	req.Header.Set("Cookie", d.Cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http_client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	str := string(content)
	// log.Println(str)
	d.isSignSuccess(str)
}

func (d JkfDriver) isSignSuccess(content string) {
	if strings.Contains(content, "成功") {
		log.Println("已经签到!", "|isSignSuccess")
		signed = true
	} else {
		log.Println("签到失败!", "|isSignSuccess")
	}
}
