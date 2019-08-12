# 基于微服务框架go-micro 封装的微服务

## 服务依赖

- 服务注册:  consul
- 消息队列:  nsq
- 链路追踪:  jaeger
- 服务间通信:  grpc
- 自动化文档:  swagger
- casbin:     权限管理
- docker-compose  集群启动
- 自定义参数验证器

依赖安装可以自己去看官方文档

##项目结构

```
├── adminApi     //管理端接口api
│   ├── docs           //swagger 文档存放位置
│   ├── handler        //接口 控制器
│   ├── middleware     //gin 中间件 包括(jwt,casbin)
│   ├── model          //model 存放目录
│   │   └── casbin     //casbin 初始化和方法存放目录
│   └── routers        //路由定义
├── api         //前端接口
│   ├── docs       //swagger 文档存放位置
│   ├── handler    //接口 控制器
│   └── middleware  //gin 中间件 包括(jwt)
├── lib         //公用库封装
│   ├── cache      //缓存 封装
│   ├── config     //配置文件 封装
│   ├── jwt        //jwt 封装
│   ├── log        //日志 (项目，gin，gorm，nsq)  封装
│   ├── orm        //orm  (gorm) 封装
│   ├── queue      //消息队列 nsq 封装
│   ├── toolfunc   //公用方法封装
│   ├── trace      //链路追踪 封装
│   └── validator  //自定义验证器
└── service      //rpc 服务项目
    ├── handler      // rpc 函数
    ├── model        // model
    ├── proto        // grpc proto生成的文件  
    └── subscriber   //消息队列 消费者监听处理
```

## 启动项目
cd 到对应的项目，go run main.go
例如启动 adminApi
```
$   cd  ./adminApi
$   go run main.go
```
