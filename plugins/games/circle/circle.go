package circle

import (
	"context"
	"fmt"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/niuhuan/mirai-bot/database/mongo"
	"github.com/niuhuan/mirai-bot/database/redis"
	"github.com/niuhuan/mirai-framework"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"strings"
	"time"
)

const id = "CIRCLE"
const name = "圈子"

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
			if content == name {
				client.ReplyText(messageInterface,
					"在群中发送以下内容 \n\n 签到 \n 积分 \n 打劫@一个人 \n ")
				return true
			}
			return false
		},
		OnGroupMessage: func(client *mirai.Client, groupMessage *message.GroupMessage) bool {
			content := client.MessageContent(groupMessage)
			if "签到" == content {
				sign(client, groupMessage)
				return true
			}
			if "积分" == content {
				points(client, groupMessage)
				return true
			}
			if strings.HasPrefix(content, "打劫") {
				rob(client, groupMessage)
				return true
			}
			return false
		},
	}
}

func points(client *mirai.Client, groupMessage *message.GroupMessage) {
	points := loadPoint(groupMessage.GroupCode, groupMessage.Sender.Uin)
	client.ReplyText(groupMessage, fmt.Sprintf("积分合计 : %v", points.Point))
}

func sign(client *mirai.Client, groupMessage *message.GroupMessage) {
	day := time.Now()
	pre := day.Add(-time.Hour * 24)
	dayStr := day.Format("2006-01-02")
	preStr := pre.Format("2006-01-02")
	lock, err := redis.TryLock(fmt.Sprintf("CIRCEL::LOCK::%v", groupMessage.GroupCode), time.Second*5, time.Second*15)
	if err != nil {
		return
	}
	defer lock.Unlock()
	last := lastSignTime(groupMessage.GroupCode, groupMessage.Sender.Uin)
	if last != nil && dayStr == last.LastDay {
		client.ReplyText(groupMessage, "您今天已经签到过")
		return
	}
	var series int
	if last != nil && last.LastDay == preStr {
		series = last.SignSeries + 1
	} else {
		series = 1
	}
	up := rand.Int()%15 + 15 // 基础积分15, 随机积分15
	up += series / 2         // 每连续签到2天多获得1积分
	saveLastSignTime(groupMessage.GroupCode, groupMessage.Sender.Uin, dayStr, series)
	points := loadPoint(groupMessage.GroupCode, groupMessage.Sender.Uin)
	incPoint(groupMessage.GroupCode, groupMessage.Sender.Uin, up)
	client.ReplyText(
		groupMessage,
		fmt.Sprintf(
			"签到成功 : \n"+
				" 连续签到 : %v 天\n"+
				" 获得积分 : %v \n"+
				" 积分合计 : %v \n\n" +
				"再接再厉, 连续签到会让获得的积分变多喔",
			series, up, points.Point+up,
		),
	)
}

func rob(client *mirai.Client, groupMessage *message.GroupMessage) {
	at := client.MessageFirstAt(groupMessage)
	if at == 0 {
		client.ReplyText(groupMessage, "您需要发送 打劫并@一个人 才能打劫他人积分")
		return
	}
	lock, err := redis.TryLock(fmt.Sprintf("CIRCEL::LOCK::%v", groupMessage.GroupCode), time.Second*5, time.Second*15)
	if err != nil {
		return
	}
	defer lock.Unlock()
	day := time.Now()
	dayStr := day.Format("2006-01-02")
	timeKey := fmt.Sprintf("CIRCEL::ROB::%v::%v::%v", dayStr, groupMessage.GroupCode, groupMessage.Sender.Uin)
	_, err = redis.GetStringError(timeKey)
	if err == nil {
		client.ReplyText(groupMessage, "每天只能打劫一次")
		return
	}
	if err == redis.Nil {
		srcPoints := loadPoint(groupMessage.GroupCode, groupMessage.Sender.Uin)
		dstPoints := loadPoint(groupMessage.GroupCode, at)
		if 30 >= dstPoints.Point {
			client.ReplyText(groupMessage, "他已经没有钱可以被打劫了")
			return
		}
		if rand.Int()%100 < 10 {
			if redis.SetString(timeKey, "1", time.Hour*24) {
				incPoint(groupMessage.GroupCode, groupMessage.Sender.Uin, 100)
				client.ReplyText(
					groupMessage,
					fmt.Sprintf("打劫时被狗咬, 丢失 %v 积分, 积分合计 : %v", 100, srcPoints.Point-100),
				)
			}
			return
		}
		if redis.SetString(timeKey, "1", time.Hour*24) {
			inc := rand.Int() % 25
			incPoint(groupMessage.GroupCode, at, -inc)
			incPoint(groupMessage.GroupCode, groupMessage.Sender.Uin, inc)
			client.ReplyText(
				groupMessage,
				fmt.Sprintf("打劫到 %v 积分, 积分合计 : %v", inc, srcPoints.Point+inc),
			)
		}
	}
}

func lastSignTime(groupCode int64, uin int64) *SignTime {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := mongo.Collection("game.circle.lastSign")
	cur, err := coll.Find(ctx, bson.M{"groupCode": groupCode, "uin": uin})
	if err != nil {
		panic(err)
	}
	defer cur.Close(ctx)
	if cur.Next(ctx) {
		var s SignTime
		err = cur.Decode(&s)
		if err != nil {
			panic(err)
		}
		return &s
	} else {
		return nil
	}
}

func saveLastSignTime(code int64, uin int64, dayStr string, series int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := mongo.Collection("game.circle.lastSign")
	_, err := coll.UpdateOne(
		ctx,
		bson.M{"uin": uin, "groupCode": code},
		bson.M{"$set": bson.M{"signSeries": series, "lastDay": dayStr}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		panic(err)
	}
}

func loadPoint(groupCode int64, uin int64) Points {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := mongo.Collection("game.circle.points")
	cur, err := coll.Find(ctx, bson.M{"groupCode": groupCode, "uin": uin})
	if err != nil {
		panic(err)
	}
	defer cur.Close(ctx)
	if cur.Next(ctx) {
		var s Points
		err = cur.Decode(&s)
		if err != nil {
			panic(err)
		}
		return s
	} else {
		return Points{
			GroupCode: groupCode,
			Uin:       uin,
			Point:     0,
		}
	}
}

func incPoint(groupCode int64, uin int64, inc int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := mongo.Collection("game.circle.points")
	_, err := coll.UpdateOne(
		ctx,
		bson.M{"uin": uin, "groupCode": groupCode},
		bson.M{"$inc": bson.M{"point": inc}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		panic(err)
	}
}

type SignTime struct {
	GroupCode  int64
	Uin        int64
	SignSeries int
	LastDay    string
}

type Points struct {
	GroupCode int64
	Uin       int64
	Point     int
}
