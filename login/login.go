package login

import (
	"bufio"
	"bytes"
	"fmt"
	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/niuhuan/mirai-bot/utils"
	"github.com/niuhuan/mirai-framework"
	logger "github.com/sirupsen/logrus"
	"github.com/tuotoo/qrcode"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var Login = false

func CmdLogin(c *mirai.Client) error {
	resp, err := c.Login()
	if err != nil {
		return err
	}
	return loginResult(c, resp)
}

func QrcodeLogin(c *mirai.Client) error {
	rsp, err := c.FetchQRCode()
	if err != nil {
		return err
	}
	fi, err := qrcode.Decode(bytes.NewReader(rsp.ImageData))
	if err != nil {
		return err
	}
	_ = os.WriteFile("qrcode.png", rsp.ImageData, 0o644)
	defer func() { _ = os.Remove("qrcode.png") }()
	logger.Infof("请使用手机QQ扫描二维码 (qrcode.png) : ")
	time.Sleep(time.Second)
	qrcodeTerminal.New().Get(fi.Content).Print()
	s, err := c.QueryQRCodeStatus(rsp.Sig)
	if err != nil {
		return err
	}
	prevState := s.State
	for {
		time.Sleep(time.Second)
		s, _ = c.QueryQRCodeStatus(rsp.Sig)
		if s == nil {
			continue
		}
		if prevState == s.State {
			continue
		}
		prevState = s.State
		switch s.State {
		case client.QRCodeCanceled:
			logger.Fatalf("扫码被用户取消.")
		case client.QRCodeTimeout:
			logger.Fatalf("二维码过期")
		case client.QRCodeWaitingForConfirm:
			logger.Infof("扫码成功, 请在手机端确认登录.")
		case client.QRCodeConfirmed:
			res, err := c.QRCodeLogin(s.LoginInfo)
			if err != nil {
				return err
			}
			return loginResult(c, res)
		case client.QRCodeImageFetch, client.QRCodeWaitingForScan:
			// ignore
		}
	}
}

func loginResult(c *mirai.Client, resp *client.LoginResponse) error {
	var err error
	console := bufio.NewReader(os.Stdin)
	for {
		if err != nil {
			logger.WithError(err).Fatal("无法登录")
			return err
		}
		var text string
		if !resp.Success {
			switch resp.Error {
			case client.SliderNeededError:
				logger.Info("请参考 https://github.com/mzdluo123/TxCaptchaHelper 获取并输入ticker")
				logger.Info("Slider url : ", resp.VerifyUrl)
				f := strings.Replace(resp.VerifyUrl, "ssl.captcha.qq.com", "txhelper.glitch.me", -1)
				logger.Info("Slider url : ", f)
				var a string
				func() {
					rsp, err := http.DefaultClient.Get(f)
					if err != nil {
						panic(err)
					}
					defer rsp.Body.Close()
					buff, err := ioutil.ReadAll(rsp.Body)
					a = string(buff)
					if err != nil {
						panic(err)
					}
				}()
				println(a)
				console.ReadString('\n')
				func() {
					rsp, err := http.DefaultClient.Get(f)
					if err != nil {
						panic(err)
					}
					defer rsp.Body.Close()
					buff, err := ioutil.ReadAll(rsp.Body)
					a = string(buff)
					if err != nil {
						panic(err)
					}
				}()
				println(a)
				resp, err = c.SubmitTicket(a)
				continue
			case client.NeedCaptcha:
				var file *os.File
				file, err = ioutil.TempFile("mirai", utils.GetSnowflakeIdString()+".jpeg")
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
	return nil
}
