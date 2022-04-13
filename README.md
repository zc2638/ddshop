# ddshop
由于上海疫情，为了提高下单成功率，让大家可以顺利吃上饭。  
此程序提供自动化抢购，若非难以存活，请给他人留一些机会！

**温馨提示:**  
1. 提前将需要购买的商品加入到购物车，并且勾上需求购买的商品
2. 在开抢前运行程序（一般6点开售，提前几分钟即可）
3. 查看日志提示，抢到后去手机上付款

**注意：** 长时间运行，可能会被封号，且行且珍惜
## 安装
### Releases
从 [Github Release](https://github.com/zc2638/ddshop/releases) 下载
### 源码
```shell
go install github.com/zc2638/ddshop/cmd/ddshop@latest
```

## 使用
1. 使用抓包工具获取 叮咚买菜上的用户 `cookie` (DDXQSESSID)
2. 使用获取到的 `cookie` 替换下面命令中的 `<custom-cookie>`
```shell
ddshop --cookie <custom-cookie>
```

自定义请求间隔(ms), 默认为 500.  
接口连续请求的时间间隔，太频繁容易被封
```shell
ddshop --cookie <custom-cookie> --interval 500
```

支持预设置支付方式  
默认支持的值：支付宝、alipay、微信、wechat
```shell
ddshop --cookie <custom-cookie> --pay-type wechat
```

Bark推送提醒 [点击查看详情](https://github.com/Finb/Bark)  
使用获取到的 `bark id` 替换下面命令中的 `<custom-bark-key>`
```shell
ddshop --cookie <custom-cookie> --bark-key <custom-bark-key>
```

## 抓包
[Charles抓包教程](https://www.jianshu.com/p/ff85b3dac157)  
微信小程序支持PC版，所以只需要安装抓包程序，打开 `叮咚买菜微信小程序`，直接进行抓包即可，无须进行手机配置。

## 声明
本项目仅供学习交流，严禁用作商业行为！  
因他人私自不正当使用造成的违法违规行为与本人无关！  
如有任何问题可联系本人删除！