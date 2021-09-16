package log

import (
	"context"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/niuhuan/mirai-bot/database/mongo"
	"github.com/niuhuan/mirai-framework/client"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	DirectionSending   = "SENDING"
	DirectionReceiving = "RECEIVING"
	TypePrivate        = "PRIVATE"
	TypeGroup          = "GROUP"
	TypeTemp           = "TEMP"
)

func NewPluginInstance() *client.Plugin {
	return &client.Plugin{
		Id: func() string {
			return "LOG"
		},
		Name: func() string {
			return "日志"
		},
		OnMessage: func(client *client.Client, messageInterface interface{}) bool {
			var m *bson.M

			if privateMessage, b := (messageInterface).(*message.PrivateMessage); b {
				buff, err := client.FormatMessageElements(privateMessage.Elements)
				if err == nil {
					m = &bson.M{
						"Direction":  DirectionReceiving,
						"Type":       TypePrivate,
						"GroupCode":  0,
						"Uin":        privateMessage.Sender.Uin,
						"Time":       privateMessage.Time,
						"InternalId": privateMessage.InternalId,
						"MsgId":      privateMessage.Id,
						"Content":    string(buff),
					}
				}
			} else if groupMessage, b := (messageInterface).(*message.GroupMessage); b {
				buff, err := client.FormatMessageElements(groupMessage.Elements)
				if err == nil {
					m = &bson.M{
						"Direction":  DirectionReceiving,
						"Type":       TypeGroup,
						"GroupCode":  groupMessage.GroupCode,
						"Uin":        groupMessage.Sender.Uin,
						"Time":       groupMessage.Time,
						"InternalId": groupMessage.InternalId,
						"MsgId":      groupMessage.Id,
						"Content":    string(buff),
					}
				}
			} else if tempMessage, b := (messageInterface).(*message.TempMessage); b {
				buff, err := client.FormatMessageElements(tempMessage.Elements)
				if err == nil {
					m = &bson.M{
						"Direction":  DirectionReceiving,
						"Type":       TypeTemp,
						"GroupCode":  tempMessage.GroupCode,
						"Uin":        tempMessage.Sender.Uin,
						"Time":       time.Now().Unix(),
						"InternalId": tempMessage.Id,
						"MsgId":      tempMessage.Id,
						"Content":    string(buff),
					}
				}
			}
			if m != nil {
				save(m)
			}
			return false
		},
	}
}

func NewListenerInstance() *client.ActionListener {
	return &client.ActionListener{
		Id: func() string {
			return "LOG"
		},
		Name: func() string {
			return "日志"
		},
		OnSendPrivateMessage: func(c *client.Client, message *message.PrivateMessage) bool {
			buff, err := c.FormatMessageElements(message.Elements)
			if err == nil {
				save(&bson.M{
					"Direction":  DirectionSending,
					"Type":       TypePrivate,
					"GroupCode":  0,
					"Uin":        message.Target,
					"Time":       message.Time,
					"InternalId": message.InternalId,
					"MsgId":      message.Id,
					"Content":    string(buff),
				})
			}
			return false
		},
		OnSendGroupMessage: func(c *client.Client, message *message.GroupMessage) bool {
			buff, err := c.FormatMessageElements(message.Elements)
			if err == nil {
				save(&bson.M{
					"Direction":  DirectionSending,
					"Type":       TypeGroup,
					"GroupCode":  message.GroupCode,
					"Uin":        0,
					"Time":       message.Time,
					"InternalId": message.InternalId,
					"MsgId":      message.Id,
					"Content":    string(buff),
				})
			}
			return false
		},
		OnSendTempMessage: func(c *client.Client, message *message.TempMessage, target int64) bool {
			buff, err := c.FormatMessageElements(message.Elements)
			if err == nil {
				save(&bson.M{
					"Direction":  DirectionSending,
					"Type":       TypeTemp,
					"GroupCode":  message.GroupCode,
					"Uin":        target,
					"Time":       time.Now().Unix(),
					"InternalId": message.Id,
					"MsgId":      message.Id,
					"Content":    string(buff),
				})
			}
			return false
		},
	}
}

func save(m *bson.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongo.Collection("log.message").InsertOne(ctx, &m)
}
