package plugins

import (
	"github.com/niuhuan/mirai-bot/plugins/games/farm"
	"github.com/niuhuan/mirai-bot/plugins/query/imglab"
	"github.com/niuhuan/mirai-bot/plugins/query/itpk"
	"github.com/niuhuan/mirai-bot/plugins/sys/ignore"
	"github.com/niuhuan/mirai-bot/plugins/sys/log"
	"github.com/niuhuan/mirai-bot/plugins/sys/menu"
	"github.com/niuhuan/mirai-bot/plugins/tools/gm"
	"github.com/niuhuan/mirai-framework"
)

func Register(c *mirai.Client) {
	// 事件监听器
	actionsListeners := []*mirai.ActionListener{
		log.NewListenerInstance(),
	}
	// 自定义组件
	cPlugins := []*mirai.Plugin{
		gm.NewPluginInstance(),
		imglab.NewPluginInstance(),
		farm.NewPluginInstance(),
		// 因为最后itpk会拦截所有私聊并做出回复, 所以一定要放在最后
		itpk.NewPluginInstance(),
	}
	// 系统组件
	sPlugins := []*mirai.Plugin{
		log.NewPluginInstance(),
		// 忽略指定用户的发言 所以放在菜单的钱main
		ignore.NewPluginInstance(),
		menu.NewPluginInstance(cPlugins),
	}
	c.SetActionListenersAndPlugins(
		actionsListeners,
		append(sPlugins, cPlugins...),
	)
}
