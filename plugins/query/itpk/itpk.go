package itpk

import (
	"encoding/json"
	"fmt"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/niuhuan/mirai-framework/client"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

const name = "闲聊"
const jokePattern = "笑话"

// NewPluginInstance 一个闲聊的插件
func NewPluginInstance() *client.Plugin {
	return &client.Plugin{
		Id: func() string {
			return "ITPK"
		},
		Name: func() string {
			return name
		},
		OnMessage: func(client *client.Client, messageInterface interface{}) bool {
			content := client.MessageContent(messageInterface)
			if strings.EqualFold(name, content) {
				printHelp(client, messageInterface)
				return true
			}
			if jokePattern == content {
				joke(client, messageInterface)
				return true
			}
			if _, ok := messageInterface.(*message.PrivateMessage); ok {
				talk(client, messageInterface)
				return true
			}
			if groupMessage, ok := messageInterface.(*message.GroupMessage); ok {
				for _, element := range groupMessage.Elements {
					if message.At == element.Type() {
						if at, ok := element.(*message.AtElement); ok {
							if at.Target == client.Uin {
								talk(client, messageInterface)
								return true
							}
						}
					}
				}
			}
			return false
		},
	}
}

func printHelp(c *client.Client, messageInterface interface{}) {
	c.ReplyText(messageInterface, "可以直接发'笑话', 或者跟我私聊, 或者@我")
}

func joke(client *client.Client, source interface{}) {
	var jockHttpRequest, _ = http.NewRequest("GET", "http://i.itpk.cn/api.php?question=笑话", nil)
	response, err := http.DefaultClient.Do(jockHttpRequest)
	if err != nil {
		logger.Error("itpk err : {}", err.Error())
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("itpk decode err : {}", err.Error())
		logger.Error("message : {}", string(body))
		return
	}
	body = body[3:] // bom 0xEF, 0xBB, 0xBF
	var jock struct {
		Title   string
		Content string
	}
	err = json.Unmarshal(body, &jock)
	if err != nil {
		logger.Error("itpk err : {}", err.Error())
		return
	}
	jock.Content = strings.ReplaceAll(jock.Content, "&nbsp;", " ")
	jock.Content = strings.ReplaceAll(jock.Content, "&amp;", ";")
	jock.Content = strings.ReplaceAll(jock.Content, "&quot;", "\"")
	jock.Content = strings.ReplaceAll(jock.Content, "&gt;", ">")
	jock.Content = strings.ReplaceAll(jock.Content, "&lt;", "<")
	jock.Content = strings.ReplaceAll(jock.Content, "<br>", "\n")
	client.ReplyText(source, fmt.Sprintf("%s\n\n%s", jock.Title, jock.Content))
}

func talk(client *client.Client, sourceMessage interface{}) {
	talkRequest, _ := http.NewRequest("GET", "http://i.itpk.cn/api.php?question="+strings.TrimSpace(client.MessageContent(sourceMessage)), nil)
	response, err := http.DefaultClient.Do(talkRequest)
	if err != nil {
		defer response.Body.Close()
	}
	body, _ := ioutil.ReadAll(response.Body)
	con := string(body)
	con = strings.ReplaceAll(con, "[cqname]", client.Nickname)
	var name string
	if groupMessage, ok := sourceMessage.(message.GroupMessage); ok {
		name = client.CardNameInGroup(groupMessage.GroupCode, groupMessage.Sender.Uin)
	} else if friend := client.FindFriend(client.MessageSenderUin(sourceMessage)); friend != nil {
		name = friend.Nickname
	} else {
		name = fmt.Sprintf("%d", client.MessageSenderUin(sourceMessage))
	}
	con = strings.ReplaceAll(con, "[name]", name)
	client.ReplyText(sourceMessage, con)
}
