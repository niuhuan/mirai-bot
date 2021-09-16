module github.com/niuhuan/mirai-bot

go 1.16

require (
	github.com/Mrs4s/MiraiGo v0.0.0-20210906051204-59288fc4dcf2
	github.com/garyburd/redigo v1.6.2
	github.com/niuhuan/mirai-framework v0.0.0
	github.com/sirupsen/logrus v1.8.1
	go.mongodb.org/mongo-driver v1.7.2
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/niuhuan/mirai-framework v0.0.0 => ./framework
