# 使用阿里云openapi更新FC函数计算

## 目录结构
```
.
├── configs/                    # 配置文件
│   └── config.json.example    # 配置模板
├── internal/                  # 私有应用代码
│   ├── config/               # 配置处理
│   └── service/              # 业务逻辑
├── pkg/                      # 公共库
└── main.go                   # 应用入口
```

## 配置
购买阿里云polarDB，按照pre环境配置购买即可，子网选择需要与函数计算在同一个VPC下，并且可以被函数计算访问
设置polarDB登陆用户和密码，并配置到config.json中(此操作已经集成到actions中)
启动应用，会自动创建数据库和表，不同环境需要手动在polarDB中插入对应数据
1. 环境配置：
```
INSERT INTO `environment_config` (environment, country_code,region) values('pre','US','us-east-1');
```
2. 阿里云账号信息：
```
INSERT INTO `aliyun_account_info` (site_client,access_key_id,access_key_secret,account_id,main_account_id,environment,region,description) 
VALUES ('client1','${access_key_id}','${access_key_secret}','${account_id}','${main_account_id}','pre','us-east-1','pre环境账号'); 
```
3. Cloudflare账号信息：
```
INSERT INTO `cloudflare_account_info` (site_client,account_id,access_key_id,access_key_secret,environment,description,country_code,api_token) 
VALUES ('client1','${account_id}','${access_key_id}','${access_key_secret}','pre','pre环境账号','US','${api_token}')
```
在运行应用之前，你需要设置你的配置：

1. 复制配置模板：
   ```bash
   cp configs/config.json.example configs/config.json
   ```

注意：永远不要将 `config.json` 文件提交到版本控制系统！

## 使用说明
```bash
1. 拉取代码
2. 运行 go mod tidy
3. 运行 go run main.go 或者 go build -o main main.go && ./main
```
### 生成swagger文档
```bash
# 安装 swag 工具
go install github.com/swaggo/swag/cmd/swag@latest
# 生成文档
swag init 
或者
swag init --parseDependency --parseInternal
```

### mac打包linux平台二进制文件
```bash
# 静态链接打包，避免 glibc 版本依赖
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o main main.go
```

## API 文档

所有需要认证的 API 都需要在请求头中携带 token：
```header
Authorization: Bearer <token>
```
### swagger
```http
GET /swagger/index.html
```

### 健康检查
#### 健康检查
```http
GET /healthcheck
```

### 认证接口
#### 登录获取 token
```http
POST /api/v1/auth/login
Content-Type: application/json

{
    "username": "admin",
    "password": "784512"
}
```

### 服务管理
#### 获取服务列表
```http
GET /api/v1/services
```

### 系统管理
#### 获取区域信息
```http
GET /api/v1/system/zones
```

### 响应格式

成功响应：
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        // 响应数据
    }
}
```

错误响应：
```json
{
    "code": 500,
    "message": "错误描述",
    "error": "详细错误信息"
}
```
