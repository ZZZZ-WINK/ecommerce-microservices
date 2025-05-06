# 电商微服务系统（Go版）

## 项目简介
这是一个基于 Go 语言开发的简单电商微服务系统，使用 gRPC 进行服务间通信。项目包含三个核心服务：
- 用户服务：处理用户注册和登录
- 商品服务：管理商品信息
- 订单服务：处理订单创建和查询

## 技术栈
- Go 1.21+
- gRPC：服务间通信
- GORM：ORM框架
- MySQL：数据存储

## 快速开始

### 1. 安装依赖
```bash
# 安装 protoc 编译器
# Windows: 下载 protoc 并添加到环境变量
# Linux/macOS: brew install protobuf

# 安装 Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 2. 生成 gRPC 代码
```bash
# Windows
.\generate.bat

# Linux/macOS
./generate.sh
```

### 3. 启动服务
```bash
# 启动用户服务
cd user-service
go run main.go

# 启动商品服务
cd product-service
go run main.go

# 启动订单服务
cd order-service
go run main.go
```

## 项目结构
```
ecommerce-microservices/
├── common/                 # 公共代码
│   └── proto/             # gRPC 协议文件
├── user-service/          # 用户服务
├── product-service/       # 商品服务
├── order-service/         # 订单服务
└── README.md             # 项目说明
```

## 开发进度
- [x] 项目基础架构
- [x] gRPC 服务定义
- [ ] 用户服务实现
- [ ] 商品服务实现
- [ ] 订单服务实现
- [ ] 数据库集成
- [ ] 服务间通信

## 项目简介
本项目为基于Go语言开发的电商微服务系统，涵盖用户、商品、订单、支付、库存、搜索等核心业务模块。系统采用主流微服务架构，服务间通过gRPC通信，支持高并发和高可用，适用于实际电商业务场景。

## 技术栈
- Gin：高性能Web框架，开发RESTful API
- GORM：ORM框架，操作MySQL数据库
- Redis：高效缓存，提升系统响应速度
- Kafka：消息队列，异步处理高并发场景
- MySQL：关系型数据库，存储核心业务数据
- Consul：服务注册与发现
- gRPC：高性能服务间通信
- Prometheus + Grafana：系统监控与可视化
- Jaeger：分布式链路追踪

## 服务说明
- **user-service**：用户注册、登录、信息管理
- **product-service**：商品信息管理、查询
- **order-service**：订单创建、查询、状态管理
- **payment-service**：支付处理、回调
- **inventory-service**：库存扣减、查询
- **search-service**：商品搜索、推荐
- **common**：公共库（gRPC proto、工具、配置等）

## 依赖环境安装
### 1. MySQL
已安装，无需重复操作。

### 2. Redis
建议使用Redis 6.x及以上版本。
- Windows可用[Memurai](https://www.memurai.com/)或[Redis官方WSL](https://redis.io/docs/getting-started/installation/install-redis-on-windows/)。
- Linux/macOS可直接`sudo apt install redis-server`或`brew install redis`。

### 3. Kafka
建议使用Kafka 2.8及以上版本。
- 下载地址：https://kafka.apache.org/downloads
- 启动命令参考README下方。

### 4. Consul
- 下载地址：https://www.consul.io/downloads
- 解压后直接运行`consul agent -dev`即可。

### 5. Prometheus & Jaeger
- Prometheus: https://prometheus.io/download/
- Jaeger: https://www.jaegertracing.io/download/

> 详细安装和启动命令见下方"本地开发环境搭建"。

## 本地开发环境搭建
### 1. 克隆项目
```bash
git clone https://github.com/你的用户名/ecommerce-microservices.git
cd ecommerce-microservices
```

### 2. 安装依赖
每个服务目录下执行：
```bash
go mod tidy
```

### 3. 启动依赖服务
#### MySQL
确保MySQL已启动，创建数据库（如`ecommerce`），并配置好账号密码。

#### Redis
```bash
# Linux/macOS
redis-server
# Windows（推荐Memurai或WSL）
memurai.exe
```

#### Kafka
```bash
# 启动Zookeeper
bin/zookeeper-server-start.sh config/zookeeper.properties
# 启动Kafka
bin/kafka-server-start.sh config/server.properties
```

#### Consul
```bash
consul agent -dev
```

#### Prometheus & Jaeger
参考各自官网文档启动。

### 4. 启动各微服务
每个服务目录下：
```bash
go run main.go
```

## 目录结构
```
user-service/         # 用户服务
product-service/      # 商品服务
order-service/        # 订单服务
payment-service/      # 支付服务
inventory-service/    # 库存服务
search-service/       # 搜索服务
common/               # 公共库（proto、工具等）
README.md             # 项目说明
```

## 联系方式
如有问题欢迎提issue或联系作者。 