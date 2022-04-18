# ddshop
购物情况有所改善，并且叮咚增加了安全风控检测，导致频繁使用程序 账号会被封禁。

该项目理论上不会再跟踪 叮咚买菜 的新更新，代码可供大家参考。  
该项目具备完整的 项目结构 以及 CICD 自动化，可供大家学习！

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
  pay_type: "wechat" # 支付方式：支付宝、alipay、微信、wechat
  channel: 3         # 通道: app => 3, 微信小程序 => 4
  interval: 100      # 连续发起请求间隔时间(ms)
  retry_count: 100   # 每次请求失败的尝试次数, -1为无限
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