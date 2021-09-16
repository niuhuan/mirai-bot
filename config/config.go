package config

import (
	"encoding/hex"
	"fmt"
	"github.com/niuhuan/mirai-bot/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
)

const configFilename = "mirai.yml"

const defaultContent = `
database:
  mongo:
    hostname: localhost
    port: 27017
    database: bot
  redis:
    hostname: localhost
    port: 6379

bot:
  account:
    uin: 123456789
    password: 5f4dcc3b5aa765d61d8327deb882cf99
`
const message = `

请修改mirai.yml
配置mongo和redis两种数据库
机器人的账号和密码(你可以使用 " echo -n password|md5 " 生成md5)

如果登录需要验证码, 您可以按照账号之前登录过的安卓手机修改device.json(可选项)

`

var Config struct {
	Database struct {
		Mongo struct {
			Hostname string
			Port     int
			Database string
		}
		Redis struct {
			Hostname string
			Port     int
		}
	}
	Bot struct {
		Account struct {
			Uin           int64
			Password      string
			PasswordBytes [16]byte
		}
	}
}

func init() {
	_, err := os.Stat(configFilename)
	if err != nil {
		err = ioutil.WriteFile(configFilename, []byte(defaultContent), os.FileMode(0644))
		utils.PanicNotNil(err)
	}
	_, err = os.Stat(configFilename)
	utils.PanicNotNil(err)
	content, _ := ioutil.ReadFile(configFilename)
	err = yaml.Unmarshal(content, &Config)
	utils.PanicNotNil(err)
	if Config.Bot.Account.Uin == 123456789 {
		fmt.Print(message)
		os.Exit(0)
	}
	regex, _ := regexp.Compile("^[0-9a-fA-F]{32}$")
	if !regex.MatchString(Config.Bot.Account.Password) {
		panic("password must be md5 32")
	}
	temp, _ := hex.DecodeString(Config.Bot.Account.Password)
	for i := 0; i < 16; i++ {
		Config.Bot.Account.PasswordBytes[i] = temp[i]
	}
}
