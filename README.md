# ddshop
由于上海疫情，为了提高下单成功率，让大家可以顺利吃上饭。  
此程序提供自动化抢购，若非难以存活，请给他人留一些机会！

## 安装
### Releases
从 [Github Release](https://github.com/zc2638/ddshop/releases) 下载
### 源码
```shell
go install github.com/zc2638/ddshop/cmd/ddshop@latest
```

## 使用
请先使用抓包工具获取 叮咚买菜上的用户 `cookie` ,然后替换下面命令中的 `<custom-cookie>`
```shell
ddshop --cookie <custom-cookie>
```

## 抓包
[Charles抓包教程](https://www.jianshu.com/p/ff85b3dac157)  
微信小程序支持PC版，所以只需要安装抓包程序，打开 `叮咚买菜微信小程序`，直接进行抓包即可，无须进行手机配置。

## 声明
本项目仅供学习交流，严禁用作商业行为！  
因他人私自不正当使用造成的违法违规行为与本人无关！  
如有任何问题可联系本人删除！