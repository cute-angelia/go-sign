package main

import (
	"log"
	"strings"
)

func main() {

	str := `<?xml version="1.0" encoding="utf-8"?>
<root><![CDATA[<script type="text/javascript" reload="1">if(typeof succeedhandle_ls=='function') {succeedhandle_ls('http://hk-bc.xyz/./', '歡迎您回來，幼兒生 a308057848，現在將轉入登錄前頁面', {'username':'a308057848','usergroup':'幼兒生','uid':'3764977','groupid':'10','syn':'0'});}hideWindow('ls');showDialog('歡迎您回來，幼兒生 a308057848，現在將轉入登錄前頁面', 'right', null, function () { window.location.href ='http://hk-bc.xyz/./'; }, 0, null, null, null, null, null, 1);</script>]]></root> iLogin
`

	if strings.Contains(str, "a308057848") {
		log.Println("登陆成功!", "|iLogin")
	} else {
		log.Println("登陆失败!", "|iLogin")
	}

}
