package farm

import (
	"context"
	"github.com/niuhuan/mirai-bot/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"time"
)

// 用户库存

func stock(groupCode int64, uin int64) Stock {
	var stock Stock
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stockCollection := mongo.Collection("game.farm.stock")
	cur, _ := stockCollection.Find(ctx, bson.M{"groupCode": groupCode, "uin": uin})
	defer cur.Close(ctx)
	if cur.Next(ctx) {
		cur.Decode(&stock)
	} else {
		stock.GroupCode = groupCode
		stock.Uin = uin
		stock.CropCount = map[string]int{}
		stockUpdate(stock)
	}
	return stock
}

func stockUpdate(stock Stock) {
	cropCount := bson.M{}
	for k, v := range stock.CropCount {
		cropCount[k] = v
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user := bson.M{"uin": stock.Uin, "groupCode": stock.GroupCode}
	update := bson.M{"$set": bson.M{"cropCount": cropCount}}
	stockCollection := mongo.Collection("game.farm.stock")
	stockCollection.UpdateOne(ctx, user, update, options.Update().SetUpsert(true))
}

type Stock struct {
	Uin       int64
	GroupCode int64
	CropCount map[string]int
}

// 用户资产

func assets(groupCode int64, uin int64) Assets {
	var assets Assets
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	assetsCollection := mongo.Collection("game.farm.assets")
	cur, _ := assetsCollection.Find(ctx, bson.M{"groupCode": groupCode, "uin": uin})
	defer cur.Close(ctx)
	if cur.Next(ctx) {
		cur.Decode(&assets)
	} else {
		assets.GroupCode = groupCode
		assets.Uin = uin
		assets.Exp = 0
		assets.Coins = 3000
		assets.Fields = 1
		assets.Ponds = 1
		assetsCollection.InsertOne(ctx, bson.M{
			"uin":       assets.Uin,
			"groupCode": assets.GroupCode,
			"exp":       assets.Exp,
			"coins":     assets.Coins,
			"fields":    assets.Fields,
			"ponds":     assets.Ponds,
		})
	}
	return assets
}

func assetsCoinsInc(code int64, uin int64, inc int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	assetsCollection := mongo.Collection("game.farm.assets")
	assetsCollection.UpdateOne(ctx, bson.M{"uin": uin, "groupCode": code}, bson.M{"$inc": bson.M{"coins": inc}})
}

func assetsFieldInc(code int64, uin int64, inc int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	assetsCollection := mongo.Collection("game.farm.assets")
	assetsCollection.UpdateOne(ctx, bson.M{"uin": uin, "groupCode": code}, bson.M{"$inc": bson.M{"fields": inc}})
}

func assetsExpInc(code int64, uin int64, up int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	assetsCollection := mongo.Collection("game.farm.assets")
	assetsCollection.UpdateOne(ctx, bson.M{"uin": uin, "groupCode": code}, bson.M{"$inc": bson.M{"exp": up}})
}

type Assets struct {
	GroupCode int64
	Uin       int64
	Exp       int64
	Coins     int64
	Fields    int
	Ponds     int
}

// 场地

func land(groupCode int64, uin int64) Land {
	var land Land
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	landCollection := mongo.Collection("game.farm.land")
	cur, _ := landCollection.Find(ctx, bson.M{"groupCode": groupCode, "uin": uin})
	defer cur.Close(ctx)
	if cur.Next(ctx) {
		cur.Decode(&land)
	} else {
		land.GroupCode = groupCode
		land.Uin = uin
		land.Fields = map[string]Field{}
		landUpdate(land)
	}
	return land
}

func landUpdate(land Land) {
	fields := bson.M{}
	for k, v := range land.Fields {
		watered := bson.M{}
		for state, uin := range v.Watered {
			watered[state] = uin
		}
		fields[k] = bson.M{
			"level":     v.Level,
			"plantTime": v.PlantTime,
			"watered":   watered,
			"stealer":   v.Stealer,
			"alerted":   v.Alerted,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user := bson.M{"groupCode": land.GroupCode, "uin": land.Uin}
	update := bson.M{"$set": bson.M{"fields": fields}}
	landCollection := mongo.Collection("game.farm.land")
	landCollection.UpdateOne(ctx, user, update, options.Update().SetUpsert(true))
}

type Field struct {
	Level     int              // 种的什么果实
	PlantTime int64            // 种植的时间
	Watered   map[string]int64 // 浇水
	Stealer   []int64          // 偷菜的人
	Alerted   []int64          // 被狗咬的人
}

type Land struct {
	GroupCode int64
	Uin       int64
	Fields    map[string]Field
}

// pets

func pets(groupCode int64, uin int64) Pets {
	var pets Pets
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	landCollection := mongo.Collection("game.farm.pets")
	cur, _ := landCollection.Find(ctx, bson.M{"groupCode": groupCode, "uin": uin})
	defer cur.Close(ctx)
	if cur.Next(ctx) {
		cur.Decode(&pets)
	} else {
		pets.GroupCode = groupCode
		pets.Uin = uin
		pets.Pets = []int{}
		petsUpdate(pets)
	}
	return pets
}

func petsUpdate(pets Pets) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user := bson.M{"groupCode": pets.GroupCode, "uin": pets.Uin}
	update := bson.M{"$set": bson.M{"pets": pets.Pets}}
	landCollection := mongo.Collection("game.farm.pets")
	landCollection.UpdateOne(ctx, user, update, options.Update().SetUpsert(true))
}

type Pets struct {
	GroupCode int64
	Uin       int64
	Pets      []int
}

// 等级/金额计算公式

func level(exp int64) int {
	for i := 21; i > 0; i-- {
		if exp >= (int64(math.Pow(float64(i), float64(4)))-1)/5 {
			return i
		}
	}
	return 0
}

func fieldPrice(currentFieldCount int) int64 {
	baseNumber := float64(currentFieldCount)
	return int64(math.Pow(2.5, .75*baseNumber)*baseNumber) * 1000
}

func cropState(crop Crop, plantTime int64, now int64) (state int, emoji string) {
	width := now - plantTime
	for i := 0; i < len(crop.StepHours); i++ {
		state = i
		emoji = crop.StepEmojis[i]
		band := 3600 * int64(crop.StepHours[i])
		if width > band {
			if i == len(crop.StepHours)-1 {
				state = MATURE
				emoji = crop.FruitEmoji
				break
			} else {
				width -= band
				continue
			}
		} else {
			break
		}
	}
	return
}

func now() int64 {
	return time.Now().Unix()
}

const MATURE = -1
