# go-nat-hole
``
服务器：启动intermediary服务
B->服务器：告诉服务器B的名称
A->服务器:我需要链接B的xx端口（服务器获取NAT-A的端口ip）
服务器->B：NAT-A:端口需要和你链接并代理到ip：端口
B:启动一个UDP服务
B->NAT-A:udp send空消息
B->服务器：udp send 说我已经准备好了，叫A来连我吧（服务器获取NAT-B的端口ip）
服务器->A:告诉NAT-B的ip和端口
A->NAT-B:开始通讯
``

## build
``
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./bootstrap/intermediary/main.go
``


## NAPT分了两个大的类型
``
Symmetric NAT型 (对称型) 需要预测端口，疯狂尝试连接才能打洞成功。暂时不实现吧
Cone NAT型（圆锥型）可以打洞

``