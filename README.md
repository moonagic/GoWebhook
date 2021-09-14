# GoWebhooks
一个使用Go编写的Github webhook接收处理服务.  
[![Go Report Card](https://goreportcard.com/badge/github.com/moonagic/GoWebhook)](https://goreportcard.com/report/github.com/moonagic/GoWebhook)  
[用Golang重写Github webhook后台服务](https://moonagic.com/make-webhook-backend-with-golang/)

## 配置范例
```yml
# vi /etc/gowebhook/config.yml
host: 127.0.0.1
port: 6606

servers:
  - identifier: dsijubhdfgbgfdnghf2
    secret: oHryZnDM8gXHfpqj
    script: /etc/gowebhook/auto_build.sh
    requestUrl: /auto_build
  - identifier: dsijubhdfgbgfdnghf
    secret: oHryZnDM8gXHfpqj1
    script: /etc/gowebhook/auto_build.sh
    requestUrl: /auto_build
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
