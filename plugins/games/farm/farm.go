package farm

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/niuhuan/mirai-bot/database/redis"
	"github.com/niuhuan/mirai-bot/utils"
	"github.com/niuhuan/mirai-framework/client"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const name = "农场"

func NewPluginInstance() *client.Plugin {
	return &client.Plugin{
		Id: func() string {
			return "FARM"
		},
		Name: func() string {
			return name
		},
		OnPrivateMessage: func(client *client.Client, privateMessage *message.PrivateMessage) bool {
			if client.MessageContent(privateMessage) == name {
				client.ReplyText(privateMessage, "农场功能只能在群中使用")
				return true
			}
			return false
		},
		OnGroupMessage: func(client *client.Client, groupMessage *message.GroupMessage) bool {
			content := client.MessageContent(groupMessage)
			if strings.EqualFold(content, "农场") {
				printMenu(client, groupMessage)
				return true
			}
			if strings.EqualFold(content, "农场帮助") {
				printHelp(client, groupMessage)
				return true
			}
			if strings.EqualFold(content, "农场商店") {
				printCrops(client, groupMessage)
				return true
			}
			if strings.EqualFold(content, "守卫商店") {
				printPets(client, groupMessage)
				return true
			}
			if strings.EqualFold(content, "购买种子") {
				printHelpBuy(client, groupMessage)
				return true
			}
			if strings.EqualFold(content, "查询种子") {
				printHelpSearch(client, groupMessage)
				return true
			}
			if strings.EqualFold(content, "购买守卫") {
				printHelpBuy(client, groupMessage)
				return true
			}
			if strings.EqualFold(content, "查询守卫") {
				printHelpSearch(client, groupMessage)
				return true
			}
			if strings.EqualFold(content, "种植") {
				printHelpPlant(client, groupMessage)
				return true
			}
			if strings.EqualFold(content, "偷菜") {
				printHelpSteal(client, groupMessage)
				return true
			}
			if strings.EqualFold(content, "我的农场") {
				printSelf(client, groupMessage)
				return true
			}
			if strings.EqualFold(content, "农场等级") {
				printLevels(client, groupMessage)
				return true
			}
			if strings.EqualFold(content, "购买土地") {
				buyField(client, groupMessage)
				return true
			}
			searchRegex, _ := regexp.Compile("^查询([\\s]+)?(\\p{Han}+)([\\s]+)?$")
			if searchRegex.MatchString(content) {
				sub := searchRegex.FindStringSubmatch(content)
				name := sub[2]
				return search(client, groupMessage, name)
			}
			buyRegex, _ := regexp.Compile("^购?买([\\s]+)?(\\p{Han}+)([\\s]+)?(\\d{1,5})?([\\s]+)?$")
			if buyRegex.MatchString(content) {
				sub := buyRegex.FindStringSubmatch(content)
				name := sub[2]
				var number int
				if len(sub[4]) > 0 {
					number, _ = strconv.Atoi(sub[4])
				} else {
					number = 1
				}
				return buy(client, groupMessage, name, number)
			}
			plantRegex, _ := regexp.Compile("^播?种植?([\\s]+)?(\\p{Han}+)([\\s]+)?([\\s]+)?$")
			if plantRegex.MatchString(content) {
				sub := plantRegex.FindStringSubmatch(content)
				name := sub[2]
				return plant(client, groupMessage, name)
			}
			if strings.EqualFold(content, "收菜") {
				collect(client, groupMessage)
				return true
			}
			if strings.Index(content, "偷菜") == 0 {
				steal(client, groupMessage)
				return true
			}
			if strings.Index(content, "浇水") == 0 {
				water(client, groupMessage)
				return true
			}
			return false
		},
	}
}

func printMenu(client *client.Client, groupMessage *message.GroupMessage) {
	client.ReplyText(groupMessage,
		" === 农场菜单 === \n\n"+
			"农场帮助\n"+
			"农场商店 守卫商店\n"+
			"购买种子 查询种子\n"+
			"购买守卫 查询守卫\n"+
			"种植 收菜 偷菜 浇水\n"+
			"我的农场 农场等级\n"+
			"购买土地 ")
}

func printHelp(client *client.Client, groupMessage *message.GroupMessage) {
	client.ReplyText(groupMessage,
		"　　农场: 机器人主人无聊开发的小游戏\n\n"+
			"货币系统: "+emojiSun+"(阳光)是农场中的基本货币\n\n"+
			"升级系统: "+emojiExp+"(经验值)可以提高农场等级\n\n"+
			"　　作物: 种植种子, 经过一段时间可以 收获"+emojiSun+"(阳光)和"+emojiExp+"(经验值)\n\n"+
			"　　土地: 土地越多, 可以同时种的种子个数\n\n"+
			"　　偷菜: 赚点小外快?\n\n"+
			"　　查询: 查询种子或者其他物品的功能 例如'查询土豆'\n\n"+
			"    守卫: 特效宠物, 防止被偷, 打盹时触发减半\n"+
			"    浇水: 获得经验值和金币, 并且增加产量, 一株植物在成熟之前每个阶段可以浇水一次")
}

func printHelpBuy(client *client.Client, groupMessage *message.GroupMessage) {
	client.ReplyText(groupMessage,
		"发送 \"购买+种子名称\" 购买相应种子, 例如 \"购买土豆\".\n\n"+
			"发送 \"购买+种子名称+数量\" 购买多个种子, 例如 \"购买土豆15\".\n\n"+
			"发送 \"购买+守卫名称\" 购买相应守卫, 例如 \"购买"+petList[0].Name+"\".\n\n"+
			"使用\"农场商店\"或者\"守卫商店\"查看列表")
}

func printHelpSearch(client *client.Client, groupMessage *message.GroupMessage) {
	client.ReplyText(groupMessage,
		"发送 \"查询+种子名称\" 查询预计收益, 例如 \"查询土豆\".\n\n"+
			"发送 \"查询+守卫名称\" 查询预计收益, 例如 \"查询"+petList[0].Name+"\".")
}

func printHelpPlant(client *client.Client, groupMessage *message.GroupMessage) {
	client.ReplyText(groupMessage,
		"发送 \"种+种子名称\" 种植作物, 例如 \"种土豆\".")
}

func printHelpSteal(client *client.Client, groupMessage *message.GroupMessage) {
	client.ReplyText(groupMessage,
		"发送 \"偷菜+@一个人\" 可以偷菜, 例如 \"偷菜@张三\".")
}

func printSelf(client *client.Client, groupMessage *message.GroupMessage) {
	assets := assets(sendUser(groupMessage))
	client.ReplyText(groupMessage, fmt.Sprintf(
		"阳光　%s　%d\n"+
			"土地　%s️　%d\n"+
			"经验　%s　%d\n"+
			"等级　%s️　%d\n",
		emojiSun, assets.Coins,
		emojiField, assets.Fields,
		emojiExp, assets.Exp, emojiLevel, level(assets.Exp),
	))
}

func printLevels(client *client.Client, groupMessage *message.GroupMessage) {
	assets := assets(sendUser(groupMessage))
	level := level(assets.Exp)
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("当前农场等级为%d级(%s%d), ", level, emojiExp, assets.Exp))
	if level >= 20 {
		builder.WriteString("您已满级.")
	} else {
		builder.WriteString(fmt.Sprintf("距离升级还需要%s%d", emojiExp, ((int64(math.Pow(float64(level+1), float64(4)))-1)/5)-assets.Exp))
	}
	client.ReplyText(groupMessage, builder.String())
}

func buyField(client *client.Client, groupMessage *message.GroupMessage) {
	// 加锁
	lock := lockUnit(sendUser(groupMessage))
	defer lock.Unlock()
	//
	assets := assets(sendUser(groupMessage))
	fieldPrice := fieldPrice(assets.Fields)
	if assets.Coins >= fieldPrice {
		assetsCoinsInc(groupMessage.GroupCode, groupMessage.Sender.Uin, -fieldPrice)
		assetsFieldInc(groupMessage.GroupCode, groupMessage.Sender.Uin, 1)
		client.ReplyText(groupMessage, "购买成功 土地+1\n"+
			fmt.Sprintf("%s ↓ %d => %d", emojiSun, fieldPrice, assets.Coins-fieldPrice))
	} else {
		client.ReplyText(groupMessage, fmt.Sprintf("购买第%d块土地需要%s%d", assets.Fields+1, emojiSun, fieldPrice))
	}
}

func search(client *client.Client, groupMessage *message.GroupMessage, name string) bool {
	for _, crop := range cropList {
		if strings.EqualFold(crop.Name, name) {
			searchCrop(client, groupMessage, crop)
			return true
		}
	}
	return false
}

func searchCrop(client *client.Client, groupMessage *message.GroupMessage, crop Crop) {
	client.ReplyText(groupMessage,
		fmt.Sprintf("%s　%s, %d级别作物, 种子售价%s%d, 成熟时间%d小时.", crop.FruitEmoji, crop.Name, crop.Level, emojiSun, crop.SeedPrice, utils.SumInts(crop.StepHours))+
			fmt.Sprintf(" 每株结出果实%d到%d枚, 预计最少收益%s%d+%s%d。",
				crop.FruitsMin, crop.FruitsMax,
				emojiSun, crop.FruitsMin*crop.FruitPrice,
				emojiExp, crop.FruitsMin*crop.FruitExp))
}

func printCrops(client *client.Client, groupMessage *message.GroupMessage) {
	// 取得数据
	assets := assets(sendUser(groupMessage))
	level := level(assets.Exp)
	stock := stock(sendUser(groupMessage))
	var builder strings.Builder
	builder.WriteString(emojiLevel + " 　　　　　　" + emojiSun + "　 " + emojiStock + "\n")
	for _, crop := range cropList {
		builder.WriteString(fmt.Sprintf("%02d　%s　%s　%d　", crop.Level, crop.FruitEmoji, crop.Name, crop.SeedPrice))
		if crop.SeedPrice < 10 {
			builder.WriteString("     ")
		} else if crop.SeedPrice < 100 {
			builder.WriteString("   ")
		} else if crop.SeedPrice < 1000 {
			builder.WriteString(" ")
		}
		if count, ok := stock.CropCount[strconv.Itoa(crop.Level)]; ok {
			builder.WriteString(fmt.Sprintf("%d", count))
		} else {
			builder.WriteString("0")
		}
		builder.WriteString("\n")
	}
	builder.WriteString(fmt.Sprintf("\n%s　%d　　　%s　%d", emojiLevel, level, emojiSun, assets.Coins))
	client.ReplyText(groupMessage, builder.String())
}

func printPets(client *client.Client, groupMessage *message.GroupMessage) {
	// 取得数据
	assets := assets(sendUser(groupMessage))
	level := level(assets.Exp)
	pets := pets(sendUser(groupMessage))
	var builder strings.Builder
	builder.WriteString(emojiLevel + " 　　　　     　　" + emojiSun + "　" + emojiStock + "\n")
	for _, pet := range petList {
		builder.WriteString(fmt.Sprintf("%02d　%s　%s　%d　", pet.Level, emojiDog, pet.Name, pet.Price))
		if utils.ContainsInt(pets.Pets, pet.Level) {
			builder.WriteString("🈶️")
		} else {
			builder.WriteString("🈚️")
		}
		builder.WriteString("\n")
	}
	builder.WriteString(fmt.Sprintf("\n%s　%d　　　%s　%d", emojiLevel, level, emojiSun, assets.Coins))
	client.ReplyText(groupMessage, builder.String())
}

func buy(client *client.Client, groupMessage *message.GroupMessage, name string, number int) bool {
	// 加锁
	lock := lockUnit(sendUser(groupMessage))
	defer lock.Unlock()
	//
	for _, crop := range cropList {
		if strings.EqualFold(crop.Name, name) {
			buyCrop(client, groupMessage, crop, number)
			return true
		}
	}
	for _, pet := range petList {
		if strings.EqualFold(pet.Name, name) {
			buyPet(client, groupMessage, pet)
			return true
		}
	}
	return false
}

func buyCrop(client *client.Client, groupMessage *message.GroupMessage, crop Crop, number int) {
	assets := assets(sendUser(groupMessage))
	level := level(assets.Exp)
	stock := stock(sendUser(groupMessage))
	if crop.Level > level {
		client.ReplyText(groupMessage, fmt.Sprintf("您不能购买超过您自身等级的作物种子, 购买%s需要%d级, 您当前为%d级. ", crop.Name, crop.Level, level))
		return
	}
	downCoin := int64(crop.SeedPrice * number)
	if downCoin > assets.Coins {
		client.ReplyText(groupMessage, fmt.Sprintf("您的阳光不足, 购买%d枚%s种子需要%d阳光, 您只有%d阳光. ", number, crop.Name, downCoin, assets.Coins))
		return
	}
	inStock, _ := stock.CropCount[strconv.Itoa(crop.Level)]
	toInStock := inStock + number
	if toInStock > 99 {
		client.ReplyText(groupMessage, "一种种子持有量不能超过99枚")
		return
	}
	stock.CropCount[strconv.Itoa(crop.Level)] = toInStock
	stockUpdate(stock)
	assetsCoinsInc(assets.GroupCode, assets.Uin, -downCoin)
	client.ReplyText(groupMessage, fmt.Sprintf("购买成功\n\n%s ↑ %d => %d\n%s ↓ %d => %d", crop.FruitEmoji, number, toInStock, emojiSun, downCoin, assets.Coins-downCoin))
}

func buyPet(client *client.Client, groupMessage *message.GroupMessage, pet Pet) {
	assets := assets(sendUser(groupMessage))
	pets := pets(sendUser(groupMessage))
	downCoin := int64(pet.Price)
	if downCoin > assets.Coins {
		client.ReplyText(groupMessage, fmt.Sprintf("您的阳光不足, 购买%s需要%d阳光, 您只有%d阳光. ", pet.Name, downCoin, assets.Coins))
		return
	}
	if utils.ContainsInt(pets.Pets, pet.Level) {
		client.ReplyText(groupMessage, "您已经有了该守卫")
		return
	}
	pets.Pets = append(pets.Pets, pet.Level)
	petsUpdate(pets)
	assetsCoinsInc(assets.GroupCode, assets.Uin, -downCoin)
	client.ReplyText(groupMessage, fmt.Sprintf("购买成功\n\n%s %s\n%s ↓ %d => %d", emojiDog, pet.Name, emojiSun, downCoin, assets.Coins-downCoin))
}

func plant(client *client.Client, groupMessage *message.GroupMessage, name string) bool {
	// 加锁
	lock := lockUnit(sendUser(groupMessage))
	defer lock.Unlock()
	//
	for _, crop := range cropList {
		if strings.EqualFold(crop.Name, name) {
			plantCrop(client, groupMessage, crop)
			return true
		}
	}
	return false
}

func plantCrop(client *client.Client, groupMessage *message.GroupMessage, crop Crop) {
	now := now()
	builder := strings.Builder{}
	assets := assets(sendUser(groupMessage))
	stock := stock(sendUser(groupMessage))
	land := land(sendUser(groupMessage))
	expUp := int64(0)
	for i := 0; i < assets.Fields; i++ {
		builder.WriteString(fmt.Sprintf("土地(%d) ", i+1))
		field, _ := land.Fields[strconv.Itoa(i)]
		if field.Level > 0 {
			cropPlanted := cropMap[field.Level]
			_, emoji := cropState(cropPlanted, field.PlantTime, now)
			builder.WriteString(fmt.Sprintf("%s (%s 已存在)", emoji, cropPlanted.Name))
		} else {
			if stock.CropCount[strconv.Itoa(crop.Level)] > 0 {
				stock.CropCount[strconv.Itoa(crop.Level)]--
				land.Fields[strconv.Itoa(i)] = Field{
					Level:     crop.Level,
					PlantTime: now,
					Watered:   map[string]int64{},
					Stealer:   []int64{},
					Alerted:   []int64{},
				}
				expUp += int64(crop.FruitExp)
				builder.WriteString(fmt.Sprintf(" => %s", crop.FruitEmoji))
			} else {
				builder.WriteString(fmt.Sprintf("%s种子不足", crop.Name))
			}
		}
		builder.WriteString("\n")
	}
	if expUp > 0 {
		stockUpdate(stock)
		landUpdate(land)
		assetsExpInc(assets.GroupCode, assets.Uin, expUp)
		builder.WriteString(fmt.Sprintf("\n%s ↑ %d => %d", emojiExp, expUp, assets.Exp+expUp))
	}
	client.ReplyText(groupMessage, builder.String())
}

func collect(client *client.Client, groupMessage *message.GroupMessage) {
	// 加锁
	lock := lockUnit(sendUser(groupMessage))
	defer lock.Unlock()
	//
	now := now()
	builder := strings.Builder{}
	assets := assets(sendUser(groupMessage))
	land := land(sendUser(groupMessage))
	//
	expUp := int64(0)
	coinsUp := int64(0)
	var waterSet []int64
	var stealerSet []int64
	for i := 0; i < assets.Fields; i++ {
		builder.WriteString(fmt.Sprintf("土地(%d) ", i+1))
		field, _ := land.Fields[strconv.Itoa(i)]
		if field.Level > 0 {
			cropPlanted := cropMap[field.Level]
			state, emoji := cropState(cropPlanted, field.PlantTime, now)
			if state == MATURE {
				var fruitNumber int
				if cropPlanted.FruitsMax > cropPlanted.FruitsMin {
					fruitNumber = (rand.Int() % (1 + cropPlanted.FruitsMax - cropPlanted.FruitsMin)) + cropPlanted.FruitsMin
				} else {
					fruitNumber = cropPlanted.FruitsMax
				}
				fruitNumber += len(field.Watered) * 1
				for _, water := range field.Watered {
					if !utils.ContainsInt64(waterSet, water) && groupMessage.Sender.Uin != water {
						waterSet = append(waterSet, water)
					}
				}
				fruitNumber -= len(field.Stealer) * 1
				builder.WriteString(fmt.Sprintf("%s (%s %d枚)", emoji, cropPlanted.Name, fruitNumber))
				if len(field.Stealer) > 0 {
					builder.WriteString(fmt.Sprintf("(被偷%d枚)", len(field.Stealer)*1))
					for _, stealer := range field.Stealer {
						if !utils.ContainsInt64(stealerSet, stealer) {
							stealerSet = append(stealerSet, stealer)
						}
					}
				}
				expUp += int64(fruitNumber * cropPlanted.FruitExp)
				coinsUp += int64(fruitNumber * cropPlanted.FruitPrice)
				delete(land.Fields, strconv.Itoa(i))
			} else {
				if _, ok := field.Watered[strconv.Itoa(state)]; ok {
					builder.WriteString(fmt.Sprintf("%s (%s 未成熟)", emoji+emojiWater, cropPlanted.Name))
				} else {
					builder.WriteString(fmt.Sprintf("%s (%s 未成熟)", emoji, cropPlanted.Name))
				}
			}
		} else {
			builder.WriteString(fmt.Sprintf("未种植"))
		}
		builder.WriteString("\n")
	}
	if expUp > 0 {
		landUpdate(land)
		assetsExpInc(assets.GroupCode, assets.Uin, expUp)
		assetsCoinsInc(assets.GroupCode, assets.Uin, coinsUp)
		if len(waterSet) > 0 {
			builder.WriteString("\n")
			builder.WriteString("帮你浇水的群友 : \n")
			for _, water := range waterSet {
				builder.WriteString("    " + client.CardNameInGroup(groupMessage.GroupCode, water) + "\n")
			}
		}
		if len(stealerSet) > 0 {
			builder.WriteString("\n")
			builder.WriteString("偷你菜的群友 : \n")
			for _, stealer := range stealerSet {
				builder.WriteString("    " + client.CardNameInGroup(groupMessage.GroupCode, stealer) + "\n")
			}
		}
		builder.WriteString(fmt.Sprintf("\n%s ↑ %d => %d\n%s ↑ %d => %d", emojiExp, expUp, assets.Exp+expUp, emojiSun, coinsUp, assets.Coins+coinsUp))
	}
	client.ReplyText(groupMessage, builder.String())
}

func steal(client *client.Client, groupMessage *message.GroupMessage) {
	// 加锁
	lock := lockUnit(sendUser(groupMessage))
	defer lock.Unlock()
	//
	firstAt := client.MessageFirstAt(groupMessage)
	if firstAt > 0 {
		builder := strings.Builder{}
		uin := groupMessage.Sender.Uin
		targetUin := firstAt
		if uin == targetUin {
			client.ReplyText(groupMessage, "你不能偷自己的菜")
			return
		} else {
			builder.WriteString("偷偷进入了 " + client.CardNameInGroup(groupMessage.GroupCode, targetUin) + "的农场\n\n")
		}
		now := now()
		targetAssets := assets(groupMessage.GroupCode, targetUin)

		// 判断被偷的人有没有狗
		targetPets := pets(groupMessage.GroupCode, targetUin)
		if utils.ContainsInt(targetPets.Pets, 1) {
			dog, _ := petMap[1]
			builder.WriteString(fmt.Sprintf("%s%s ", emojiDog, dog.Name))
			sleepy := (int(math.Abs(float64(rand.Int63()-targetUin))) % 100) < 50
			alertPercentage := 20
			if sleepy {
				builder.WriteString("正在瞌睡 ")
				alertPercentage /= 2
			}
			alert := (int(math.Abs(float64(rand.Int63()-targetUin))) % 100) < alertPercentage
			if alert {
				assetsCoinsInc(groupMessage.GroupCode, groupMessage.Sender.Uin, -100)
				builder.WriteString("把你咬了 损失 " + emojiSun + "100")
				client.ReplyText(groupMessage, builder.String())
				return
			}
		}

		targetLand := land(groupMessage.GroupCode, targetUin)
		expUp := int64(0)
		coinsUp := int64(0)
		for i := 0; i < targetAssets.Fields; i++ {
			builder.WriteString(fmt.Sprintf("土地(%d) ", i+1))
			field, _ := targetLand.Fields[strconv.Itoa(i)]
			if field.Level > 0 {
				cropPlanted := cropMap[field.Level]
				state, emoji := cropState(cropPlanted, field.PlantTime, now)
				if state != MATURE {
					builder.WriteString(fmt.Sprintf("%s (%s 未成熟)", emoji, cropPlanted.Name))
				} else {
					if utils.ContainsInt64(field.Stealer, groupMessage.Sender.Uin) {
						builder.WriteString(fmt.Sprintf("%s (%s 偷过了)", emoji, cropPlanted.Name))
					} else if len(field.Stealer) >= 2 {
						builder.WriteString(fmt.Sprintf("%s (%s 快被偷光了)", emoji, cropPlanted.Name))
					} else {
						expUp += int64(1 * cropPlanted.FruitExp)
						coinsUp += int64(1 * cropPlanted.FruitPrice)
						field.Stealer = append(field.Stealer, uin)
						targetLand.Fields[strconv.Itoa(i)] = field
						builder.WriteString(fmt.Sprintf("%s (%s %d枚)", emoji, cropPlanted.Name, 1))
					}
				}
			} else {
				builder.WriteString(fmt.Sprintf("未种植"))
			}
			builder.WriteString("\n")
		}
		if expUp > 0 {
			assets := assets(sendUser(groupMessage))
			landUpdate(targetLand)
			assetsExpInc(groupMessage.GroupCode, uin, expUp)
			assetsCoinsInc(groupMessage.GroupCode, uin, coinsUp)
			builder.WriteString(fmt.Sprintf("\n%s ↑ %d => %d\n%s ↑ %d => %d", emojiExp, expUp, assets.Exp+expUp, emojiSun, coinsUp, assets.Coins+coinsUp))
		}
		client.ReplyText(groupMessage, builder.String())
	} else {
		printHelpSteal(client, groupMessage)
	}
}

func water(client *client.Client, groupMessage *message.GroupMessage) {
	// 加锁
	lock := lockUnit(sendUser(groupMessage))
	defer lock.Unlock()
	//
	now := now()
	uin := groupMessage.Sender.Uin
	targetUin := uin
	for _, element := range groupMessage.Elements {
		if element.Type() == message.At {
			if at, ok := element.(*message.AtElement); ok {
				targetUin = at.Target
			}
			break
		}
	}
	targetAssets := assets(groupMessage.GroupCode, targetUin)
	targetLand := land(groupMessage.GroupCode, targetUin)
	builder := strings.Builder{}
	if uin != targetUin {
		builder.WriteString(client.CardNameInGroup(groupMessage.GroupCode, targetUin) + "的农场\n\n")
	} else {
		builder.WriteString("浇水@一个人可以为群友浇水\n\n")
	}
	expUp := int64(0)
	for i := 0; i < targetAssets.Fields; i++ {
		builder.WriteString(fmt.Sprintf("土地(%d) ", i+1))
		field, _ := targetLand.Fields[strconv.Itoa(i)]
		if field.Level > 0 {
			cropPlanted := cropMap[field.Level]
			state, emoji := cropState(cropPlanted, field.PlantTime, now)
			if state == MATURE {
				builder.WriteString(fmt.Sprintf("%s (%s 已成熟)", emoji, cropPlanted.Name))
			} else {
				if _, ok := field.Watered[strconv.Itoa(state)]; ok {
					builder.WriteString(fmt.Sprintf("%s (%s 无需浇水)", emoji+emojiWater, cropPlanted.Name))
				} else {
					targetLand.Fields[strconv.Itoa(i)].Watered[strconv.Itoa(state)] = uin
					expUp += int64(cropPlanted.FruitExp)
					builder.WriteString(fmt.Sprintf("%s (%s 浇水成功)", emoji+emojiRain, cropPlanted.Name))
				}
			}
		} else {
			builder.WriteString(fmt.Sprintf("未种植"))
		}
		builder.WriteString("\n")
	}
	if expUp > 0 {
		assets := assets(sendUser(groupMessage))
		landUpdate(targetLand)
		assetsExpInc(groupMessage.GroupCode, uin, expUp)
		builder.WriteString(fmt.Sprintf("\n%s ↑ %d => %d", emojiExp, expUp, assets.Exp+expUp))
	}
	client.ReplyText(groupMessage, builder.String())
}

func sendUser(groupMessage *message.GroupMessage) (groupCode int64, uin int64) {
	return groupMessage.GroupCode, groupMessage.Sender.Uin
}

func lockUnit(groupCode int64, uin int64) *redis.Lock {
	return redis.SaaSLock(fmt.Sprintf("BOT::GAME::FARM::%v::%v::LOCK", groupCode, uin), time.Minute)
}
