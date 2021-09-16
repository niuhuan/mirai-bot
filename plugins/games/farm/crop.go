package farm

import (
	"encoding/json"
)

var (
	cropList []Crop
	cropMap  = map[int]Crop{}
)

func init() {
	json.Unmarshal([]byte(cropJson), &cropList)
	for index := 0; index < len(cropList); index++ {
		crop := cropList[index]
		cropMap[crop.Level] = crop
	}
}

type Crop struct {
	Level      int
	Name       string
	SeedPrice  int
	FruitsMin  int
	FruitsMax  int
	FruitPrice int
	FruitExp   int
	StepHours  []int
	StepEmojis []string
	FruitEmoji string
}

const cropJson = `
[
  {
    "level": 1,
    "name": "土豆",
    "seedPrice": 10,
    "fruitsMin": 8,
    "fruitsMax": 12,
    "fruitPrice": 4,
    "fruitExp": 4,
    "stepHours": [
      1,
      2,
      3
    ],
    "stepEmojis": [
      "\uD83C\uDF31",
      "\uD83C\uDF31",
      "\uD83C\uDF8D"
    ],
    "fruitEmoji": "\uD83E\uDD54"
  },
  {
    "level": 2,
    "name": "萝卜",
    "seedPrice": 20,
    "fruitsMin": 10,
    "fruitsMax": 15,
    "fruitPrice": 8,
    "fruitExp": 4,
    "stepHours": [
      1,
      2,
      3
    ],
    "stepEmojis": [
      "\uD83C\uDF31",
      "\uD83C\uDF8D",
      "\uD83C\uDF8D"
    ],
    "fruitEmoji": "\uD83E\uDD55"
  },
  {
    "level": 3,
    "name": "花生",
    "seedPrice": 30,
    "fruitsMin": 15,
    "fruitsMax": 17,
    "fruitPrice": 8,
    "fruitExp": 4,
    "stepHours": [
      1,
      3,
      4
    ],
    "stepEmojis": [
      "\uD83C\uDF31",
      "\uD83C\uDF8D",
      "\uD83C\uDF3F"
    ],
    "fruitEmoji": "\uD83E\uDD5C"
  },
  {
    "level": 4,
    "name": "番茄",
    "seedPrice": 40,
    "fruitsMin": 10,
    "fruitsMax": 15,
    "fruitPrice": 20,
    "fruitExp": 9,
    "stepHours": [
      1,
      3,
      4
    ],
    "stepEmojis": [
      "\uD83C\uDF31",
      "\uD83C\uDF8D",
      "\uD83C\uDF3F"
    ],
    "fruitEmoji": "\uD83C\uDF45"
  },
  {
    "level": 5,
    "name": "茄子",
    "seedPrice": 50,
    "fruitsMin": 10,
    "fruitsMax": 15,
    "fruitPrice": 25,
    "fruitExp": 12,
    "stepHours": [
      2,
      4,
      5
    ],
    "stepEmojis": [
      "\uD83C\uDF31",
      "\uD83C\uDF8D",
      "\uD83C\uDF3F"
    ],
    "fruitEmoji": "\uD83C\uDF46"
  },
  {
    "level": 6,
    "name": "辣椒",
    "seedPrice": 120,
    "fruitsMin": 20,
    "fruitsMax": 25,
    "fruitPrice": 25,
    "fruitExp": 12,
    "stepHours": [
      2,
      4,
      5
    ],
    "stepEmojis": [
      "\uD83C\uDF31",
      "\uD83C\uDF8D",
      "\uD83C\uDF3E"
    ],
    "fruitEmoji": "\uD83C\uDF36"
  },
  {
    "level": 7,
    "name": "蘑菇",
    "seedPrice": 140,
    "fruitsMin": 25,
    "fruitsMax": 30,
    "fruitPrice": 25,
    "fruitExp": 12,
    "stepHours": [
      2,
      4,
      6
    ],
    "stepEmojis": [
      "\uD83C\uDF31",
      "\uD83C\uDF8D",
      "\uD83C\uDF3E"
    ],
    "fruitEmoji": "\uD83C\uDF44"
  },
  {
    "level": 8,
    "name": "玉米",
    "seedPrice": 160,
    "fruitsMin": 30,
    "fruitsMax": 35,
    "fruitPrice": 50,
    "fruitExp": 20,
    "stepHours": [
      2,
      4,
      6
    ],
    "stepEmojis": [
      "\uD83C\uDF31",
      "\uD83C\uDF8D",
      "\uD83C\uDF3E"
    ],
    "fruitEmoji": "\uD83C\uDF3D"
  },
  {
    "level": 11,
    "name": "苹果",
    "seedPrice": 220,
    "fruitsMin": 30,
    "fruitsMax": 35,
    "fruitPrice": 60,
    "fruitExp": 30,
    "stepHours": [
      3,
      6,
      8
    ],
    "stepEmojis": [
      "\uD83C\uDF31",
      "\uD83C\uDF8D",
      "\uD83C\uDF33"
    ],
    "fruitEmoji": "\uD83C\uDF4E"
  },
  {
    "level": 13,
    "name": "雪梨",
    "seedPrice": 260,
    "fruitsMin": 30,
    "fruitsMax": 35,
    "fruitPrice": 70,
    "fruitExp": 30,
    "stepHours": [
      3,
      6,
      8
    ],
    "stepEmojis": [
      "\uD83C\uDF31",
      "\uD83C\uDF8D",
      "\uD83C\uDF33"
    ],
    "fruitEmoji": "\uD83C\uDF50"
  },
  {
    "level": 15,
    "name": "桃子",
    "seedPrice": 300,
    "fruitsMin": 30,
    "fruitsMax": 35,
    "fruitPrice": 100,
    "fruitExp": 70,
    "stepHours": [
      3,
      6,
      8
    ],
    "stepEmojis": [
      "\uD83C\uDF31",
      "\uD83C\uDF8D",
      "\uD83C\uDF33"
    ],
    "fruitEmoji": "\uD83C\uDF51"
  },
  {
    "level": 17,
    "name": "橙子",
    "seedPrice": 510,
    "fruitsMin": 30,
    "fruitsMax": 35,
    "fruitPrice": 150,
    "fruitExp": 100,
    "stepHours": [
      3,
      6,
      8
    ],
    "stepEmojis": [
      "\uD83C\uDF31",
      "\uD83C\uDF8D",
      "\uD83C\uDF33"
    ],
    "fruitEmoji": "\uD83C\uDF4A"
  },
  {
    "level": 19,
    "name": "柠檬",
    "seedPrice": 999,
    "fruitsMin": 30,
    "fruitsMax": 35,
    "fruitPrice": 200,
    "fruitExp": 150,
    "stepHours": [
      3,
      6,
      8
    ],
    "stepEmojis": [
      "\uD83C\uDF31",
      "\uD83C\uDF8D",
      "\uD83C\uDF33"
    ],
    "fruitEmoji": "\uD83C\uDF4B"
  }
]
`
