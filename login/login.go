package login

import (
	"bufio"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/niuhuan/mirai-bot/utils"
	nc "github.com/niuhuan/mirai-framework/client"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

func CmdLogin(c *nc.Client) {
	console := bufio.NewReader(os.Stdin)
	resp, err := c.Login()
	for {
		if err != nil {
			logger.WithError(err).Fatal("unable to login")
		}
		var text string
		if !resp.Success {
			switch resp.Error {
			case client.SliderNeededError:
				if client.SystemDeviceInfo.Protocol == client.AndroidPhone {
					logger.Warn("Android手机协议不支持滑动验证")
					logger.Warn("请使用其他客户端类型")
					os.Exit(2)
				}
				c.AllowSlider = false
				c.Disconnect()
				resp, err = c.Login()
				continue
			case client.NeedCaptcha:
				file, err := ioutil.TempFile("mirai", utils.GetSnowflakeIdString()+".jpeg")
				func() {
					utils.GetSnowflakeIdString()
					utils.PanicNotNil(err)
					defer file.Close()
					file.Write(resp.CaptchaImage)
				}()
				fmt.Print("请输入验证码 : 图片位置 : " + file.Name())
				text, _ := console.ReadString('\n')
				resp, err = c.SubmitCaptcha(strings.ReplaceAll(text, "\n", ""), resp.CaptchaSign)
				continue
			case client.SMSNeededError:
				fmt.Println("QQ开启了设备锁, 需要发送短信,  输入YES进行发送")
				fmt.Printf("短信发送到 %s ? (yes)", resp.SMSPhone)
				t, _ := console.ReadString('\n')
				t = strings.TrimSpace(t)
				if t != "yes" {
					os.Exit(2)
				}
				if !c.RequestSMS() {
					logger.Warnf("无法获取短信验证码")
					os.Exit(2)
				}
				logger.Warn("请输入短信验证码: ")
				text, _ = console.ReadString('\n')
				resp, err = c.SubmitSMS(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\r", ""))
				continue
			case client.SMSOrVerifyNeededError:
				fmt.Println("开启了设备锁:")
				fmt.Println("1. 发送短信到 ", resp.SMSPhone)
				fmt.Println("2. 扫描二维码")
				fmt.Print("输入 (1,2):")
				text, _ = console.ReadString('\n')
				text = strings.TrimSpace(text)
				switch text {
				case "1":
					if !c.RequestSMS() {
						logger.Warnf("无法获取短信验证码")
						os.Exit(2)
					}
					logger.Warn("请输入短信验证码: ")
					text, _ = console.ReadString('\n')
					resp, err = c.SubmitSMS(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\r", ""))
					continue
				case "2":
					fmt.Printf("设备锁 -> %v\n", resp.VerifyUrl)
					os.Exit(2)
				default:
					fmt.Println("不正确的输入")
					os.Exit(2)
				}
			case client.UnsafeDeviceError:
				fmt.Printf("设备锁 -> %v\n", resp.VerifyUrl)
				os.Exit(2)
			case client.OtherLoginError, client.UnknownLoginError:
				logger.Fatalf("登录失败: %v", resp.ErrorMessage)
				os.Exit(3)
			}
		}
		break
	}
	logger.Info("登录成功, 加载通讯录...")
	c.ReloadFriendList()
	c.ReloadGroupList()
	logger.Info("加载完成")
}
