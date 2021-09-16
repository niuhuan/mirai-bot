package ignore

import (
	"github.com/niuhuan/mirai-bot/utils"
	"github.com/niuhuan/mirai-framework/client"
)

var ignoreUidArray = []int64{
	2854196310, // Q群管家
	2854196312, // 表情包老铁
	2854196306, // 微软小冰
}

func NewPluginInstance() *client.Plugin {
	return &client.Plugin{
		Id: func() string {
			return "IGNORE"
		},
		Name: func() string {
			return "忽略"
		},
		OnMessage: func(client *client.Client, messageInterface interface{}) bool {
			if utils.ContainsInt64(ignoreUidArray, client.MessageSenderUin(messageInterface)) {
				return true
			}
			return false
		},
	}
}
