package main

import (
	"github.com/niuhuan/mirai-bot/config"
	"github.com/niuhuan/mirai-bot/database/mongo"
	"github.com/niuhuan/mirai-bot/database/redis"
	"github.com/niuhuan/mirai-bot/login"
	"github.com/niuhuan/mirai-bot/plugins"
	"github.com/niuhuan/mirai-framework"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
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
	// client := mirai.NewClientMd5(config.Config.Bot.Account.Uin, config.Config.Bot.Account.PasswordBytes)
	client := mirai.NewClient(0, "")
	// 注册插件
	plugins.Register(client)
	// 登录
	buff, err := ioutil.ReadFile("session.token")
	if err == nil {
		err = client.TokenLogin(buff)
	}
	if err != nil {
		err = login.QrcodeLogin(client)
	}
	// login.CmdLogin(client)
	if err == nil{
		ioutil.WriteFile("session.token", client.GenToken(), os.FileMode(0600))
		logger.Info("登录成功, 加载通讯录...")
		client.ReloadFriendList()
		client.ReloadGroupList()
		logger.Info("加载完成")
		login.Login = true
		// 等待退出信号
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt, os.Kill)
		<-ch
	}
}
