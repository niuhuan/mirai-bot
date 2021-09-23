package gm

import (
	client2 "github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/niuhuan/mirai-framework"
	"regexp"
	"strconv"
	"strings"
)

const id = "GROUP_MANAGER"
const name = "群管"

var banRegexp, _ = regexp.Compile("^(\\s+)?b(\\s+)?([0-9]{1,5})(\\s+)?([smhd]?)(\\s+)?$")

func mana(qqClient *mirai.Client, groupMessage *message.GroupMessage) bool {
	groupInfo := qqClient.FindGroup(groupMessage.GroupCode)
	if groupInfo != nil {
		senderInfo := groupInfo.FindMember(groupMessage.Sender.Uin)
		botInfo := groupInfo.FindMember(qqClient.Uin)
		if senderInfo != nil && botInfo != nil {
			if senderInfo.Permission == client2.Member {
				qqClient.ReplyText(groupMessage, "您必须是管理员才能使用管理指令")
			} else if botInfo.Permission == client2.Member {
				qqClient.ReplyText(groupMessage, "机器人必须是管理员才能使用管理指令")
			} else {
				return true
			}
		}
	}
	return false
}

func NewPluginInstance() *mirai.Plugin {
	return &mirai.Plugin{
		Id: func() string {
			return id
		},
		Name: func() string {
			return name
		},
		OnPrivateMessage: func(client *mirai.Client, privateMessage *message.PrivateMessage) bool {
			if client.MessageContent(privateMessage) == name {
				client.ReplyText(privateMessage, "群管功能只能在群中使用")
				return true
			}
			return false
		},
		OnGroupMessage: func(client *mirai.Client, groupMessage *message.GroupMessage) bool {
			elements := groupMessage.Elements
			if text, ok := (elements[0]).(*message.TextElement); ok {
				if banRegexp.MatchString(text.Content) {
					if mana(client, groupMessage) {
						matches := banRegexp.FindStringSubmatch(text.Content)
						source, _ := strconv.Atoi(matches[3])
						switch matches[5] {
						case "m":
							source *= 60
						case "h":
							source *= 60 * 60
						case "d":
							source *= 60 * 60 * 24
						}
						if source > 60*60*24*29 {
							client.ReplyText(groupMessage, "最多禁言29天")
						} else {
							mLen := 0
							groupInfo := client.FindGroup(groupMessage.GroupCode)
							if groupInfo != nil {
								for _, element := range elements {
									if at, ok := element.(*message.AtElement); ok {
										target := groupInfo.FindMember(at.Target)
										if target != nil && target.Manageable() {
											target.Mute(uint32(source))
											mLen++
										}
									}
								}
							}
							if mLen > 0 {
								client.ReplyText(groupMessage, "操作成功")
							}
						}
					}
					return true
				}
			}
			return false
		},
		OnMessage: func(client *mirai.Client, messageInterface interface{}) bool {
			content := client.MessageContent(messageInterface)
			if strings.EqualFold(name, content) {
				client.ReplyText(messageInterface,
					"群管理功能\n"+
						"当机器人是管理员, 且发送者为管理员时可生效\n\n"+
						" 批量禁言",
				)
				return true
			}
			if strings.EqualFold("批量禁言", content) {
				client.ReplyText(messageInterface,
					"b+禁言时间 @一个或多个人\n\n"+
						"比如禁言张三12小时 : b12h @张三 \n\n"+
						"比如禁言张三李四12天 : b12h @张三 @李四 \n\n"+
						" s 秒, m 分, h 小时, d 天\n\n"+
						"b0 则解除禁言",
				)
				return true
			}
			return false
		},
	}
}
