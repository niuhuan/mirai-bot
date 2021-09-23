package menu

import (
	"fmt"
	"github.com/niuhuan/mirai-framework"
	"strings"
)

func NewPluginInstance(customerPlugins []*mirai.Plugin) *mirai.Plugin {
	return &mirai.Plugin{
		Id: func() string {
			return "MENU"
		},
		Name: func() string {
			return "菜单"
		},
		OnMessage: func(client *mirai.Client, messageInterface interface{}) bool {
			content := client.MessageContent(messageInterface)
			if strings.EqualFold("菜单", content) {
				builder := strings.Builder{}
				builder.WriteString("菜单 : ")
				for i := 0; i < len(customerPlugins); i++ {
					builder.WriteString(fmt.Sprintf("\n♦️ %s", (*customerPlugins[i]).Name()))
				}
				client.ReplyText(messageInterface, builder.String())
				return true
			}
			return false
		},
	}
}
