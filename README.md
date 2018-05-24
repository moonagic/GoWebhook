# GoWebhooks
一个使用Go编写的Github webhook接收处理服务.  
[用Golang重写Github webhook后台服务](https://moonagic.com/make-webhook-backend-with-golang/)

## 配置
```
# vi /etc/gowebhook/config
{
  "requestUrl": "在Webhooks中设定的Payload URL",
  "secret": "在Webhooks中设定的secret",
  "host": "监听地址",
  "port": "监听端口",
  "script": "用于部署的shell脚本"
}
```

## 简单的systemd启动服务
```
[Unit]
Description=gowebhook service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/gowebhook

[Install]
WantedBy=multi-user.target
```
