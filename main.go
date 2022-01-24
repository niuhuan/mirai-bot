package main

import (
	"github.com/niuhuan/mirai-bot/config"
	"github.com/niuhuan/mirai-bot/database/mongo"
	"github.com/niuhuan/mirai-bot/database/redis"
	"github.com/niuhuan/mirai-bot/login"
	"github.com/niuhuan/mirai-bot/plugins"
	"github.com/niuhuan/mirai-framework"
	"os"
	"os/signal"
)

func main() {
	// 检查DeviceInfo, Config 初始化配置
	config.InitDeviceInfo()
	config.InitConfig()
	mongo.InitMongo()
	redis.InitRedis()
	// 校验数据库链接是否正常
	mongo.Test()
	redis.Test()
	// 新建客户端
	client := mirai.NewClientMd5(config.Config.Bot.Account.Uin, config.Config.Bot.Account.PasswordBytes)
	// 注册插件
	plugins.Register(client)
	// 登录
	login.CmdLogin(client)
	// 等待退出信号
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}
