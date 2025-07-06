# sing-box-sub-converter

sing-box-sub-converter 是一个用于合并和转换 sing-box 配置的工具，它提供了 HTTP API 接口，可以从多个订阅源获取节点信息，并将其合并到自定义的配置模板中。

## 项目架构

### 核心组件

1. **服务器 (Server)**
   - 提供 HTTP API 接口
   - 处理配置生成请求
   - 支持快速启动模式
   - 静态文件服务（提供Web界面）

2. **转换器 (Converter)**
   - 支持多种订阅格式解析（Clash、普通文本等）
   - 将不同格式的节点转换为 sing-box 格式
   - 处理节点标签、前缀和表情符号
   - 自动处理重复节点名称

3. **模板系统 (Template)**
   - 加载配置模板
   - 将节点合并到模板中
   - 支持节点过滤功能（包含/排除关键词）
   - 支持按订阅标签分组节点

4. **订阅获取 (Fetcher)**
   - 从远程URL获取订阅内容
   - 从本地文件获取订阅内容
   - 支持自定义 User-Agent
   - 支持安全目录限制（对本地文件）

5. **配置管理 (Config)**
   - 管理全局配置和订阅信息
   - 支持多个订阅源
   - 配置文件热重载（自动检测配置变更）
   - 提供默认配置生成

### 数据流

1. 客户端发送请求到 API 接口
2. 服务器加载配置模板
3. 获取订阅内容（远程URL或本地文件）
4. 尝试不同解析器解析订阅内容为节点列表
5. 处理节点（添加前缀、表情符号等）
6. 应用过滤规则（如果模板中定义）
7. 将节点合并到配置模板中
8. 返回最终配置给客户端

## 如何运行

### Docker 方式

```shell
docker run -d -p 5000:5000 mystery0/sing-box-sub-converter:latest
```

### 环境变量

| 环境变量 | 描述 | 默认值 |
|---------|------|-------|
| SERVER_PORT | 服务器监听端口 | 5000 |
| TEMPLATE_DIR | 配置模板目录 | config_templates |
| SUB_CONFIG_HOME | 订阅配置文件目录 | 当前工作目录 |
| SAFE_DIR | 本地文件订阅的安全目录 | 无（不限制） |

## API 接口

### 1. 生成配置

**接口**: `/api/generate`

**方法**: GET

**请求参数**:

| 参数名 | 类型   | 必填 | 描述                   |
|--------|--------|------|------------------------|
| file   | string | 是   | 配置模板文件名（不含扩展名） |

**返回数据**:

成功时返回 JSON 格式的 sing-box 配置。

**示例请求**:
```
GET /api/generate?file=fileName
```

**错误响应**:

| 状态码 | 响应                                  | 描述               |
|--------|---------------------------------------|-------------------|
| 400    | {"error": "Missing template file parameter"} | 缺少模板文件参数     |
| 400    | {"error": "Failed to load template"} | 加载模板失败         |
| 400    | {"error": "Failed to process subscribes"} | 处理订阅失败         |
| 400    | {"error": "Failed to merge config"} | 合并配置失败         |

### 2. 快速启动

**接口**: `/api/quickstart/*url`

**方法**: GET

**请求参数**:

| 参数名 | 类型   | 必填 | 描述                   |
|--------|--------|------|------------------------|
| url    | string | 是   | 订阅 URL（作为路径的一部分） |
| file   | string | 是   | 配置模板文件名（不含扩展名） |

**返回数据**:

成功时返回 JSON 格式的 sing-box 配置。

**示例请求**:
```
GET /api/quickstart/https://example.com/sub?file=openwrt
```

**错误响应**:

| 状态码 | 响应                                  | 描述               |
|--------|---------------------------------------|-------------------|
| 400    | {"error": "Missing subscription URL"} | 缺少订阅 URL        |
| 400    | {"error": "Missing template file parameter"} | 缺少模板文件参数     |
| 400    | {"error": "Failed to load template"} | 加载模板失败         |
| 400    | {"error": "Failed to process subscribes"} | 处理订阅失败         |
| 400    | {"error": "Failed to merge config"} | 合并配置失败         |

## 配置说明

### 订阅配置

在 `providers.json` 文件中配置订阅源：

```json
{
  "subscribes": [
    {
      "url": "订阅URL，支持本地文件 file://文件绝对路径",
      "tag": "订阅标签",
      "prefix": "节点前缀",
      "userAgent": "自定义User-Agent"
    }
  ],
  "prefix": true,
  "emoji": true,
  "exclude_protocol": "",
  "showSubInNodes": true
}
```

**配置选项说明：**

- `subscribes`: 订阅源列表
- `prefix`: 是否在节点名称前添加前缀
- `emoji`: 是否在节点名称前添加国家/地区表情符号
- `exclude_protocol`: 排除指定协议的节点
- `showSubInNodes`: 是否显示订阅信息（剩余流量和到期时间）作为节点。当设置为 `true` 时，系统会解析订阅响应头中的 `subscription-userinfo` 信息，并创建包含剩余流量和到期天数的节点。

### 配置模板

配置模板存放在 `config_template` 目录下，可通过 TEMPLATE_DIR 环境变量修改读取目录，使用 JSON 格式。模板中可以使用以下占位符：

- `{all}`: 表示所有节点
- `{tag名}`: 表示特定标签的节点

模板还支持过滤器功能，可以根据关键词包含或排除节点：

```json
{
  "filter": [
    {
      "action": "include",
      "keywords": ["关键词1", "关键词2"],
      "for": "订阅标签"
    },
    {
      "action": "exclude",
      "keywords": ["关键词3", "关键词4"],
      "for": "订阅标签"
    }
  ]
}
```

## 支持的订阅格式

### Clash 格式
支持解析 Clash 配置文件中的代理节点，包括以下协议：
- Shadowsocks
- Trojan
- VLESS
- Hysteria2

### 普通文本格式
支持解析包含以下格式的文本文件：
- URI 格式的 Shadowsocks 链接 (ss://)
- 其他协议的链接（根据前缀自动识别）

## 模板示例

### 基本模板
```json
{
  "outbounds": [
    {
      "type": "selector",
      "tag": "proxy",
      "outbounds": ["{all}"]
    }
  ]
}
```

### 分组模板
```json
{
  "outbounds": [
    {
      "type": "selector",
      "tag": "proxy",
      "outbounds": ["auto", "manual"]
    },
    {
      "type": "urltest",
      "tag": "auto",
      "outbounds": ["{all}"]
    },
    {
      "type": "selector",
      "tag": "manual",
      "outbounds": ["{all}"]
    }
  ]
}
```

### 带过滤器的模板
```json
{
  "outbounds": [
    {
      "type": "selector",
      "tag": "proxy",
      "outbounds": ["hk", "jp", "sg", "us"]
    },
    {
      "type": "selector",
      "tag": "hk",
      "filter": [
        {
          "action": "include",
          "keywords": ["香港", "HK", "Hong Kong"],
          "for": "all"
        }
      ],
      "outbounds": ["{all}"]
    },
    {
      "type": "selector",
      "tag": "jp",
      "filter": [
        {
          "action": "include",
          "keywords": ["日本", "JP", "Japan"],
          "for": "all"
        }
      ],
      "outbounds": ["{all}"]
    }
  ]
}
```

## Web 界面

sing-box-sub-converter 提供了一个简单的 Web 界面，可通过浏览器访问：

```
http://localhost:5000/
```

通过 Web 界面，您可以：
- 查看可用的配置模板
- 生成配置文件
- 测试订阅链接

## 使用方法

1. 配置 `providers.json` 文件，添加订阅源
2. 创建或修改配置模板
3. 启动服务器（默认端口 5000，可通过 SERVER_PORT 环境变量修改）
4. 通过 API 接口或 Web 界面获取生成的配置
5. 将生成的配置保存为 JSON 文件并加载到 sing-box 中使用
