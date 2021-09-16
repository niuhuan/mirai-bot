package config

import (
	"github.com/Mrs4s/MiraiGo/client"
	"io/ioutil"
	"os"
)

const deviceFilePath = "device.json"

func InitDeviceInfo() {
	_, err := os.Stat(deviceFilePath)
	if err == nil {
		buff, _ := ioutil.ReadFile(deviceFilePath)
		client.SystemDeviceInfo.ReadJson(buff)
	} else {
		client.GenRandomDevice()
		ioutil.WriteFile(deviceFilePath, client.SystemDeviceInfo.ToJson(), os.FileMode(0644))
	}
}
