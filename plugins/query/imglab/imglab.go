package imglab

import (
	"crypto/tls"
	"github.com/niuhuan/mirai-framework"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
)

var httpClient = &http.Client{Transport: &http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
}}

const id = "IMG_LIB"
const name = "图库"

func NewPluginInstance() *mirai.Plugin {
	return &mirai.Plugin{
		Id: func() string {
			return id
		},
		Name: func() string {
			return name
		},
		OnMessage: func(client *mirai.Client, messageInterface interface{}) bool {
			content := client.MessageContent(messageInterface)
			if content == "图库" {
				client.ReplyText(messageInterface,
					"发送以下内容 \n\n 动漫壁纸 \n 美女壁纸 \n 风景壁纸 \n 手机动漫壁纸 \n 手机美女壁纸 \n 手机风景壁纸 \n 随机插画 \n 随机老婆")
				return true
			}
			if content == "动漫壁纸" {
				reply(client, messageInterface, dmRequest)
				return true
			}
			if content == "美女壁纸" {
				reply(client, messageInterface, meiziRequest)
				return true
			}
			if content == "风景壁纸" {
				reply(client, messageInterface, fengjingRequest)
				return true
			}
			if content == "手机动漫壁纸" {
				reply(client, messageInterface, dmSjRequest)
				return true
			}
			if content == "手机美女壁纸" {
				reply(client, messageInterface, meiziSjRequest)
				return true
			}
			if content == "手机风景壁纸" {
				reply(client, messageInterface, fengjingSjRequest)
				return true
			}
			if content == "随机插画" {
				reply(client, messageInterface, illustrationRequest)
				return true
			}
			if content == "随机老婆" {
				reg, _ := regexp.Compile("src=\"//(.+\\.jpg)\"")
				do, err := httpClient.Do(wfRequest)
				if err == nil {
					defer do.Body.Close()
					bo, err := ioutil.ReadAll(do.Body)
					if err == nil {
						str := string(bo)
						find := reg.FindStringSubmatch(str)
						if find != nil {
							req, _ := http.NewRequest("GET", "https://"+find[1], nil)
							reply(client, messageInterface, req)
						}
					}
				}
				return true
			}
			return false
		},
	}
}


var dmRequest, _ = http.NewRequest("GET", "http://api.molure.cn/sjbz/api.php?lx=dongman", nil)
var dmSjRequest, _ = http.NewRequest("GET", "http://api.molure.cn/sjbz/api.php?method=mobile&lx=dongman", nil)
var meiziRequest, _ = http.NewRequest("GET", "http://api.molure.cn/sjbz/api.php?lx=meizi", nil)
var meiziSjRequest, _ = http.NewRequest("GET", "http://api.molure.cn/sjbz/api.php?method=mobile&lx=meizi", nil)
var fengjingRequest, _ = http.NewRequest("GET", "http://api.molure.cn/sjbz/api.php?lx=fengjing", nil)
var fengjingSjRequest, _ = http.NewRequest("GET", "http://api.molure.cn/sjbz/api.php?method=mobile&lx=fengjing", nil)

var illustrationRequest, _ = http.NewRequest("GET", "https://api.mz-moe.cn/img.php", nil)
var wfRequest, _ = http.NewRequest("GET", "https://img.xjh.me/random_img.php", nil)

func reply(client *mirai.Client, messageInterface interface{}, request *http.Request) {
	do, err := httpClient.Do(request)
	if err == nil {
		defer do.Body.Close()
		bo, err := ioutil.ReadAll(do.Body)
		if err == nil {
			img, err := client.UploadReplyImage(messageInterface, bo)
			if err == nil {
				client.ReplyRawMessage(messageInterface, client.MakeReplySendingMessage(messageInterface).Append(img))
			} else {
				logger.Error("IMG UPLOAD ERROR : {}", err.Error())
			}
		} else {
			logger.Error("IMG READ ERROR : {}", err.Error())
		}
	} else {
		logger.Error("IMG REQUEST ERROR : {}", err.Error())
	}
}

