package farm

import "encoding/json"

var (
	petList []Pet
	petMap  = map[int]Pet{}
)

func init() {
	json.Unmarshal([]byte(petJson), &petList)
	for index := 0; index < len(petList); index++ {
		pet := petList[index]
		petMap[pet.Level] = pet
	}
}

type Pet struct {
	Level int    `json:"level"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

const petJson = `
[
{
"level": 1,
"name": "斗牛犬",
"price": 10000
}
]
`
