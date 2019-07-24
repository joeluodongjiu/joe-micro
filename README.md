# 基于微服务框架go-micro 封装的微服务脚手架

## 服务依赖

- 服务注册:  consul
- 消息队列:  nsq
- 链路追踪:  jaeger
- 服务间通信:  grpc
- 自动化文档:  swagger

依赖安装可以自己去看官方文档

##项目结构

```
├── api          //接口项目
│   ├── docs         //swagger 文档存放位置
│   └── handler      //接口控制器
├── lib          
│   ├── cache        //缓存 封装
│   ├── config       //配置文件 封装
│   ├── log          //日志 (项目，gin，gorm，nsq)  封装
│   ├── orm          // orm  (gorm) 封装
│   ├── queue        // 消息队列 nsq 封装
│   └── trace        // 链路追踪 封装
└── service      //rpc 服务项目
    ├── handler      // rpc 函数
    ├── model        // model
    ├── proto        // grpc proto生成的文件  
    └── subscriber   //消息队列 消费者监听处理

```

##
启动rpc服务
```
go run ./service/main.go
```
启动api服务
```
go run ./api/main.go
```

##  TODO

- casbin          权限管理
- docker-compose  集群启动
