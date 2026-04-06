# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 常用命令

```bash
# 构建
go build ./...

# 运行
go run main.go

# 运行所有测试
go test ./...

# 运行单个测试文件
go test ./converter/ -run Test_processDuplicateNodeTag -v

# 运行单个测试包
go test ./utils/ -v

# 更新依赖
go mod tidy
```

## 环境变量

| 变量名 | 说明 | 默认值 |
|---|---|---|
| `DEBUG` | 启用调试日志 | 无 |
| `SERVER_PORT` | 监听端口 | `5000` |
| `TEMPLATE_DIR` | 配置模板目录 | `config_templates` |
| `SUB_CONFIG_HOME` | 订阅配置文件目录 | 当前工作目录 |
| `SAFE_DIR` | 本地文件订阅的安全目录限制 | 无 |
| `SUB_URL` | 未挂载 providers.json 时的默认订阅链接 | `http://127.0.0.1:1080/test.txt` |

## 架构概览

### 数据流

```
HTTP 请求 → server/server.go
  → template.GetConfigTemplate()       // 加载 JSON 模板文件
  → converter.ProcessSubscribes()      // 获取并解析订阅
      → fetcher.FetchSubscription()    // 拉取内容（本地文件或远程 URL）
      → parser.Parse()                 // 尝试各解析器（Clash → 纯文本）
      → addEmoji/addPrefix/processDuplicateNodeTag  // 节点后处理
  → template.MergeToConfig()           // 将节点填充到模板占位符
  → 返回最终 sing-box JSON 配置
```

### 核心包职责

- **`config/`** — 全局配置（`providers.json`），使用 `viper` + `fsnotify` 实现热重载；`types.go` 定义 `ProvidersGlobalConfig` 和 `Subscription`
- **`fetcher/`** — 订阅内容获取，两种实现：`file.go`（`file://` 前缀）和 `remote.go`（HTTP/HTTPS），同时解析 `Subscription-Userinfo` 响应头为 `SubInfo`
- **`converter/`** — 订阅解析与节点处理的核心；`converter.go` 定义 `SubscriptionParser` 接口；`parsers.go` 定义 `ProcessSubscribes` 主流程；`singbox.go` 负责 Shadowsocks 转 sing-box 格式的兜底实现
  - **`converter/clash/`** — 解析 Clash YAML 格式，支持 SS/Trojan/VLESS/Hysteria2/VMess
  - **`converter/content/`** — 解析纯文本 URI 格式（目前仅支持 `ss://`）
  - **`converter/types/`** — 共享类型 `ProxyNode` 和 `ProxyNodeType` 常量
- **`template/`** — 加载配置模板（`load.go`）、替换 `{all}` / `{tag名}` 占位符（`template.go`）、应用 `filter` 规则（`filter.go`）
- **`utils/`** — 国家/地区匹配（`countrymatcher.go`）、表情符号处理（`emoji.go`）、编码工具（`encoding.go`）
- **`server/`** — Gin HTTP 服务器，路由：`GET /api/generate`、`GET /api/quickstart/*url`、`GET /`（内嵌 Web 界面）
- **`api/`** — Vercel 无服务器函数入口（`vercel.go`），通过 POST body 接受模板 JSON

### 添加新协议支持

1. 在 `converter/clash/` 中新增解析文件（参考 `shadowsocks.go`）
2. 在 `converter/clash/clash.go` 的 `Parse` 方法中注册新协议
3. 在 `converter/types/types.go` 中添加 `ProxyNodeType` 常量
4. 实现 `Convert2SingBoxOutbounds` 方法将节点转为 sing-box 出站格式
5. 若需支持纯文本 URI 格式，在 `converter/content/` 中扩展

### 模板占位符规则

- `{all}` — 替换为当前 outbound 过滤后的全部节点
- `{tag名}` — 替换为来自指定订阅 tag 的节点（`ProxyNode.FromSub` 字段）
- `filter` 字段（定义在单个 outbound 内）— 处理完后会从输出中删除，不会出现在最终配置里

## 代码规范

- 所有注释、日志消息、错误文本均使用中文
- 日志使用标准库 `log/slog`（JSON 格式输出）
- 配置通过环境变量注入，不使用命令行参数
