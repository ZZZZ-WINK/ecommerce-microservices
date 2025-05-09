# 电商微服务系统（Go版）

## 项目简介
这是一个基于 Go 语言开发的电商微服务系统，采用微服务架构，使用 gRPC 进行服务间通信。项目包含三个核心服务：
- 用户服务：处理用户注册、登录和信息管理
- 商品服务：管理商品信息，支持缓存和搜索
- 订单服务：处理订单创建、查询和状态管理

## 技术栈
- Go 1.21+
- gRPC：服务间通信
- GORM：ORM框架
- MySQL：数据存储
- Redis：缓存和分布式锁
- JWT：用户认证

## 已实现功能

### 用户服务（user-service:50051）
- ✓ 用户注册
- ✓ 用户登录（JWT认证）
- ✓ 获取用户信息
- ✓ 更新用户信息
- ✓ 删除用户
- ✓ 用户列表（分页）

### 商品服务（product-service:50052）
- ✓ 创建商品
- ✓ 获取商品详情（Redis缓存）
- ✓ 商品列表（分页，关键词搜索）
- ✓ 更新商品
- ✓ 删除商品
- ✓ 首页商品列表缓存

### 订单服务（order-service:50053）
- ✓ 创建订单（分布式锁防并发）
- ✓ 获取订单详情
- ✓ 订单列表（分页）
- ✓ 更新订单状态
- ✓ 删除订单
- ✓ 库存检查

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
# 生成所有服务的代码
.\generate.bat

# 生成指定服务的代码
.\generate.bat -service user     # 只生成用户服务的代码
.\generate.bat -service product  # 只生成商品服务的代码
.\generate.bat -service order    # 只生成订单服务的代码

# 清理并重新生成代码
.\generate.bat -clean           # 清理所有生成的代码
.\generate.bat -service user -clean  # 清理并重新生成用户服务的代码
```

### 3. 配置数据库
确保 MySQL 已启动，创建数据库：
```sql
CREATE DATABASE ecommerce CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 启动服务
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
│   ├── middleware/        # 中间件（JWT认证等）
│   └── proto/            # gRPC 协议文件
├── user-service/         # 用户服务
│   ├── model/           # 数据模型
│   └── service/         # 业务逻辑
├── product-service/      # 商品服务
│   ├── model/           # 数据模型
│   └── service/         # 业务逻辑
├── order-service/        # 订单服务
│   ├── model/           # 数据模型
│   └── service/         # 业务逻辑
└── README.md            # 项目说明
```

## 技术特性
- ✓ 微服务架构
- ✓ gRPC 服务间通信
- ✓ JWT 用户认证
- ✓ Redis 缓存
- ✓ Redis 分布式锁
- ✓ 数据库事务
- ✓ 分页查询
- ○ 服务注册与发现
- ○ 配置中心
- ○ 链路追踪
- ○ 监控系统
- ○ 日志系统
- ○ 单元测试
- ○ CI/CD
- ○ 容器化部署
- ○ 限流和熔断
- ○ 消息队列

## 开发环境
- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- Windows 11 + WSL2

## 本地开发
1. 确保 MySQL 和 Redis 服务已启动
2. 克隆项目并进入项目目录
3. 安装依赖：`go mod tidy`
4. 生成 gRPC 代码
5. 启动各个服务

## 待优化项
1. 添加服务注册与发现（Consul/etcd）
2. 集成配置中心（Nacos/Apollo）
3. 添加链路追踪（Jaeger）
4. 集成监控系统（Prometheus + Grafana）
5. 完善日志系统
6. 添加单元测试
7. 配置 CI/CD
8. 容器化部署
9. 添加限流和熔断
10. 集成消息队列（Kafka）

## 贡献指南
欢迎提交 Issue 和 Pull Request。

## 许可证
MIT License 