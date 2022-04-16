# ddshop
由于上海疫情，为了提高下单成功率，让大家可以顺利吃上饭。  
此程序提供自动化抢购，若非难以存活，请给他人留一些机会！

**温馨提示:**  
1. 提前将需要购买的商品加入到购物车，并且勾上需求购买的商品
2. 在开抢前运行程序（一般6点开售，提前几秒即可）
3. 查看日志提示，抢到后去手机上付款

**注意：** 长时间运行，可能会被封号，且行且珍惜
## 安装
### Releases
从 [Github Release](https://github.com/zc2638/ddshop/releases) 下载
### Docker
```shell
docker pull zc2638/ddshop:latest
```
### 源码
```shell
go install github.com/zc2638/ddshop/cmd/ddshop@latest
```

## 使用

1. 创建配置文件`config.yaml`
2. 修改配置文件中的配置项
3. 执行程序

### 配置
点击查看 [完整配置](./config/config.yaml)
```yaml
bark:
  key: ""               # Bark 通知推送的 Key
push_plus:
  token: ""             # Push Plus 通知推送的 Token
  
# 自动任务的配置，不配置 periods 将持续执行
regular:
  success_interval: 100 # 执行成功 再次执行的间隔时间(ms), -1为停止继续执行
  fail_interval: 100    # 执行失败 再次执行的间隔时间(ms), -1为停止继续执行
  periods: # 执行周期
    - start: "05:59"
      end: "06:10"
    - start: "08:29"
      end: "08:35"

# 叮咚买菜的配置
ddmc:
  cookie: ""         # 使用抓包工具获取 叮咚买菜上的用户 `cookie` (DDXQSESSID)
  interval: 100      # 连续发起请求间隔时间(ms)
  payType: "wechat"  # 支付方式：支付宝、alipay、微信、wechat
  channel: 3         # 通道: app => 3, 微信小程序 => 4
```
### 命令行工具
执行以下命令前，将 `<custom-config-path>` 替换为实际的配置文件路径，例如：`config/config.yaml`
```shell
ddshop -c <custom-config-path>
```
### Docker
执行以下命令前，将 `<custom-config-dir>` 替换成宿主机存放配置文件的目录  
```shell
docker run --name ddshop -it -v <custom-config-dir>:/work/config zc2638/ddshop 
```

## 抓包
[Charles抓包教程](https://www.jianshu.com/p/ff85b3dac157)  
[Charles 抓包 PC端微信小程序](https://blog.csdn.net/z2181745/article/details/123002569)  
微信小程序支持PC版，所以只需要安装抓包程序，打开 `叮咚买菜微信小程序`，直接进行抓包即可，无须进行手机配置。

## 声明
本项目仅供学习交流，严禁用作商业行为！  
因他人私自不正当使用造成的违法违规行为与本人无关！  
如有任何问题可联系本人删除！