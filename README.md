mirai-bot
=====
一个基于MiraiFramework的QQ机器人


# 设计思路

[MiraiFramework](https://github.com/niuhuan/mirai-framework)

所有的功能都是由插件完成, 事件发生时, 调度器对插件循环调用, 插件响应是否处理该事件, 直至有插件响应事件, 插件发生异常, 或插件轮训结束, 最后日志结果被记录, 事件响应周期结束。
![img.png](images/invoke.png)

# 运行须知

- 第一次运行 会生成 mirai.yml 和 device.json, 修改后启动即可
- 第一次登录 您可以安装安卓软件DeviceInfo, 参照内容修改device.json, 并将protocol改为2(安卓手表)/1(安卓手机)将绕过设备锁
- 以后运行将很少失败, 您可以使用docker启动
- 本bot使用了redis和mongo, mongo和redis解压可直接使用, 如果您没有条件下载, 可以删除database包 / 需要数据库的plugin 之后启动.

依赖关系
- 农场插件 : mongo / redis

# 功能展示

![](images/plugin01.jpg)
![](images/plugin02.jpg)
![](images/plugin03.jpg)
![](images/plugin04.jpg)
![](images/plugin05.jpg)
![](images/plugin06.jpg)