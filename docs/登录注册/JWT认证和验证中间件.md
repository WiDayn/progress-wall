# JWT 认证 和 验证中间件
## 概述

后端使用JWT (JSON Web Token) 进行身份认证。前端需要通过登录接口获取 `accessToken`，然后在后续需要认证的API请求中携带此token。

## 认证流程

1. 新用户通过注册接口创建账号（已有账号可跳过此步骤）
2. 前端调用登录接口获取 `accessToken`
3. 前端存储 `accessToken`（推荐使用 localStorage 或 sessionStorage）
4. 后续需要认证的API请求在请求头中携带 `Authorization: Bearer <accessToken>`
5. 后端验证token，验证失败返回 401 Unauthorized

## API接口

### 登录接口

**接口地址**: `POST /api/auth/login`

**请求头**:
```
Content-Type: application/json
```

**请求体**:
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**说明**:
- `username`: 支持用户名或邮箱地址
- `password`: 用户密码

**成功响应** (200 OK):
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "nickname": "系统管理员",
    "avatar": "",
    "phone": "",
    "status": 1,
    "last_login": "2024-01-01T12:00:00Z",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

**错误响应**:

1. 请求参数错误 (400 Bad Request):
```json
{
  "error": "请求参数错误"
}
```

2. 用户名或密码错误 (401 Unauthorized):
```json
{
  "error": "用户名或密码错误"
}
```

3. 账户已被禁用 (403 Forbidden):
```json
{
  "error": "账户已被禁用"
}
```

4. 服务器错误 (500 Internal Server Error):
```json
{
  "error": "生成token失败"
}
```

### 注册接口

**接口地址**: `POST /api/auth/register`

**请求头**:
```
Content-Type: application/json
```

**请求体**:
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "nickname": "测试用户"
}
```

**说明**:
- `username`: 用户名，必填，唯一标识
- `email`: 邮箱地址，必填，必须是有效的邮箱格式，唯一标识
- `password`: 密码，必填，至少6位字符
- `nickname`: 昵称，可选

**成功响应** (201 Created):
```json
{
  "user": {
    "id": 2,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户",
    "avatar": "",
    "phone": "",
    "status": 1,
    "last_login": null,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

**错误响应**:

1. 请求参数错误 (400 Bad Request):
```json
{
  "error": "请求参数错误"
}
```

2. 密码格式不正确 (400 Bad Request):
```json
{
  "error": "密码格式不正确"
}
```

3. 用户名或邮箱已存在 (409 Conflict):
```json
{
  "error": "用户名或邮箱已存在"
}
```

4. 服务器错误 (500 Internal Server Error):
```json
{
  "error": "创建用户失败"
}
```

**注意**: 注册成功后，用户需要调用登录接口获取 `accessToken` 才能进行后续的认证操作。

## Token使用

### 在请求中使用Token

对于需要认证的API请求，必须在请求头中添加 `Authorization` 字段：

**请求头格式**:
```
Authorization: Bearer <accessToken>
```

**说明**:
- `<accessToken>` 为登录接口返回的 `accessToken` 值
- 格式必须为 `Bearer <token>`，Bearer 和 token 之间有一个空格

### 后端中间件使用

后端使用 `AuthMiddleware` 中间件来验证token。验证成功后，中间件会将用户信息存储到请求上下文中，后端路由处理函数可以通过以下方式获取：

**Gin框架获取方式**:
```go
userID := c.GetUint("user_id")      // 获取用户ID (uint类型)
username := c.GetString("username")  // 获取用户名 (string类型)
```

**说明**:
- `user_id`: 当前登录用户的ID（uint类型）
- `username`: 当前登录用户的用户名（string类型）
- 这些信息在token验证成功后自动设置到请求上下文
- 在受保护的路由中，可以直接使用这些信息进行业务处理

## 错误处理

### 401 Unauthorized 响应

当token验证失败时，后端会返回 401 状态码，可能的原因和响应：

1. **未提供token**:
```json
{
  "error": "未提供认证token"
}
```

2. **Token格式错误**:
```json
{
  "error": "认证token格式错误"
}
```

3. **Token无效或过期**:
```json
{
  "error": "无效的token"
}
```

## Token过期

JWT token具有过期时间，默认配置为24小时（可通过环境变量 `JWT_EXPIRE_HOURS` 配置）。

当token过期时，后端会返回 401 Unauthorized 状态码，前端需要引导用户重新登录。

## 注意事项

1. **Token安全**:
   - 不要在URL参数中传递token
   - 不要将token存储在可被XSS攻击访问的地方（如全局变量）

2. **Token存储**:
   - `localStorage` / `sessionStorage`
   - 根据安全需求选择存储方式

3. **错误处理**:
   - 统一处理401错误，自动跳转登录

4. **Token刷新**:
   - 当前未实现token刷新机制
   - token过期后需要用户重新登录

## 测试账号 / 初始账号（在初始化数据库时会自动创建）

- **用户名**: `admin`
- **密码**: `admin123`
- **邮箱**: `admin@example.com`

